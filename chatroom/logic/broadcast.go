package logic

import (
	"log"

	"github.com/summerKK/go-code-snippet-library/chatroom/global"
)

type broadcaster struct {
	// 聊天室所有用户
	users map[string]*User

	// 进入聊天室用户
	enteringChannel chan *User
	// 离开聊天室用户
	leavingChannel chan *User
	// 消息广播
	messageChannel chan *Message

	// 判断用户名是否重复.从`checkCanInChannel`返回结果
	checkUserChannel  chan string
	checkCanInChannel chan bool

	// 获取用户列表.从`usersChannel`返回结果
	requestUsersChannel chan struct{}
	usersChannel        chan []*User
}

var Broadcaster = &broadcaster{
	users: make(map[string]*User),

	enteringChannel: make(chan *User),
	leavingChannel:  make(chan *User),

	messageChannel: make(chan *Message, global.MessageQueueLen),

	checkUserChannel:  make(chan string),
	checkCanInChannel: make(chan bool),

	requestUsersChannel: make(chan struct{}),
	usersChannel:        make(chan []*User),
}

func (b *broadcaster) Start() {
	for {
		select {
		case u := <-b.enteringChannel:
			// 用户进去聊天室.把用户放在用户列表中
			b.users[u.Nickname] = u
			// 给用户发送离线消息
			OfflineProcessor.Send(u)
		case u := <-b.leavingChannel:
			// 用户离开聊天室.需要关闭结束消息的goroutine.避免内存泄露
			delete(b.users, u.Nickname)
			// 关闭消息channel
			close(u.MessageChannel)
		case msg := <-b.messageChannel:
			for _, u := range b.users {
				// 不用发送给自己
				if u.Uid == msg.User.Uid {
					continue
				}

				u.MessageChannel <- msg
			}
			// 保存最近消息
			OfflineProcessor.Save(msg)

		case nickname := <-b.checkUserChannel:
			_, ok := b.users[nickname]
			b.checkCanInChannel <- !ok
		case <-b.requestUsersChannel:
			var userList []*User
			// append是增加元素.make的时候已经初始化了容量.所以这里不能指定容量
			// ❌ userList := make([]*User, len(b.users))
			for _, user := range b.users {
				userList = append(userList, user)
			}
			b.usersChannel <- userList
		}
	}
}

func (b *broadcaster) CanEnterRoom(nickname string) bool {
	b.checkUserChannel <- nickname

	return <-b.checkCanInChannel
}

func (b *broadcaster) Broadcast(msg *Message) {
	if len(b.messageChannel) >= global.MessageQueueLen {
		log.Println("消息发送队列已满")
	}
	b.messageChannel <- msg
}

func (b *broadcaster) UserEntering(user *User) {
	b.enteringChannel <- user
}

func (b *broadcaster) UserLeaving(user *User) {
	b.leavingChannel <- user
}

func (b *broadcaster) UserList() []*User {
	b.requestUsersChannel <- struct{}{}

	return <-b.usersChannel
}

func (b *broadcaster) CheckUserOnline(nickname string) bool {
	_, ok := b.users[nickname]

	return ok
}
