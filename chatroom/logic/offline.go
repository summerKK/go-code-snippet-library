package logic

import (
	"container/ring"

	"github.com/spf13/viper"
)

type offlineProcessor struct {
	n int
	// 保存所有用户的离线消息(n 条)
	recentRing *ring.Ring
	// 保存某个用户的离线消息(n 条)
	userRing map[string]*ring.Ring
}

var OfflineProcessor = newOfflineProcessor()

func newOfflineProcessor() *offlineProcessor {
	n := viper.GetInt("offline-num")

	return &offlineProcessor{
		n:          n,
		recentRing: ring.New(n),
		userRing:   make(map[string]*ring.Ring),
	}
}

// 保存离线消息
func (o *offlineProcessor) Save(msg *Message) {
	if msg.Type != MsgTypeNormal {
		return
	}

	// 这里主要解决了 `slice`的地址是指针数组的问题.通过make.生成一个新的slice.然后把元素再copy回去
	// 这样copyMsg和msg的 Ats指向的内存地址就不一样了
	copyMsg := *msg
	copyMsg.Ats = make([]string, len(msg.Ats))
	copy(copyMsg.Ats, msg.Ats)
	copyMsg.Ats = nil

	o.recentRing.Value = &copyMsg
	// ring 是一个环形链表,这里指针下移
	o.recentRing = o.recentRing.Next()

	for _, nickname := range msg.Ats {
		nickname = nickname[1:]
		var (
			r  *ring.Ring
			ok bool
		)
		// 如果用户在线就不需要重复艾特他了
		if Broadcaster.CheckUserOnline(nickname) {
			continue
		}

		if r, ok = o.userRing[nickname]; !ok {
			r = ring.New(o.n)
		}

		r.Value = msg
		o.userRing[nickname] = r.Next()
	}
}

// 用户登录发送历史消息
func (o *offlineProcessor) Send(user *User) {
	// 遍历离线消息列表
	o.recentRing.Do(func(value interface{}) {
		if value != nil {
			user.MessageChannel <- value.(*Message)
		}
	})

	if user.isNew {
		return
	}

	// 艾特(@)消息
	if r, ok := o.userRing[user.Nickname]; ok {
		r.Do(func(value interface{}) {
			if value != nil {
				user.MessageChannel <- value.(*Message)
			}
		})

		delete(o.userRing, user.Nickname)
	}
}
