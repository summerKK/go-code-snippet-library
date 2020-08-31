package logic

import (
	"context"
	"errors"
	"io"
	"regexp"
	"sync/atomic"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type User struct {
	Uid            int           `json:"uid"`
	Nickname       string        `json:"nickname"`
	EnterAt        time.Time     `json:"enter_at"`
	Addr           string        `json:"addr"`
	MessageChannel chan *Message `json:"-"`
	Token          string        `json:"token"`

	conn *websocket.Conn

	isNew bool
}

var globalUid uint64 = 0

// 系统用户
var System = &User{}

// 解析@
var regAtUsers = regexp.MustCompile(`@[^\s@]{2,20}`)

func NewUser(nickname, token, addr string, conn *websocket.Conn) *User {
	user := &User{
		Uid:            int(atomic.AddUint64(&globalUid, 1)),
		Nickname:       nickname,
		EnterAt:        time.Now(),
		Addr:           addr,
		MessageChannel: make(chan *Message, 32),
		Token:          token,
		conn:           conn,
	}

	// 已经登录过的用户
	//if user.Token != "" {
	//	uid, err := parseTokenAndValidate(token, nickname)
	//	if err == nil {
	//		user.Uid = uid
	//	}
	//}
	//
	//if user.Uid == 0 {
	//	user.Uid = int(atomic.AddUint32(&globalUid, 1))
	//	user.Token = genToken(user.Uid, user.Nickname)
	//	user.isNew = true
	//}

	return user
}

// 给用户发送消息
func (u *User) SendMessage(ctx context.Context) {
	for msg := range u.MessageChannel {
		wsjson.Write(ctx, u.conn, msg)
	}
}

// 接收用户的发送消息,并发送到聊天室
func (u *User) ReceiveMessage(ctx context.Context) error {
	var (
		receiveMsg map[string]string
		err        error
	)

	for {
		err = wsjson.Read(ctx, u.conn, &receiveMsg)
		if err != nil {
			var closeErr websocket.CloseError
			// 如果错误是正常关闭则忽略
			if errors.As(err, &closeErr) {
				return nil
			}
			if errors.Is(err, io.EOF) {
				return nil
			}

			return err
		}

		// 内容发送到聊天室
		sendMsg := NewMsg(u, receiveMsg["content"], receiveMsg["client_time"])

		// 解析content,看看@了谁
		sendMsg.Ats = regAtUsers.FindAllString(sendMsg.Content, -1)

		Broadcaster.Broadcast(sendMsg)
	}
}
