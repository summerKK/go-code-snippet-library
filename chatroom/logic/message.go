package logic

import (
	"fmt"
	"time"

	"github.com/spf13/cast"
)

type messageType int

const (
	// 普通 用户消息
	MsgTypeNormal messageType = iota
	// 当前用户欢迎消息
	MsgTypeWelcome
	// 用户进入
	MsgTypeEnter
	// 用户离开
	MsgTypeLeave
	// 错误消息
	MsgTypeError
)

type Message struct {
	// 哪个用户发送消息
	User    *User       `json:"user"`
	Type    messageType `json:"type"`
	Content string      `json:"content"`
	MsgTime time.Time   `json:"msg_time"`

	ClientSendTime time.Time `json:"client_send_time"`

	// @了谁
	Ats []string `json:"ats"`
}

func NewErrorMsg(content string) *Message {
	return &Message{
		User:    System,
		Type:    MsgTypeError,
		Content: content,
		MsgTime: time.Now(),
	}
}

func NewWelcomeMsg(user *User) *Message {
	return &Message{
		User:    user,
		Type:    MsgTypeWelcome,
		Content: fmt.Sprintf("欢迎: %s 进入聊天室", user.Nickname),
		MsgTime: time.Now(),
	}
}

func NewEnterMsg(user *User) *Message {
	return &Message{
		User:    user,
		Type:    MsgTypeEnter,
		Content: fmt.Sprintf("%s 进入聊天室", user.Nickname),
		MsgTime: time.Now(),
	}
}

func NewLeavingMsg(user *User) *Message {
	return &Message{
		User:    user,
		Type:    MsgTypeLeave,
		Content: fmt.Sprintf("%s 离开了聊天室", user.Nickname),
		MsgTime: time.Now(),
	}
}

func NewMsg(user *User, content, clientSendTime string) *Message {
	msg := &Message{
		User:    user,
		Type:    MsgTypeNormal,
		Content: content,
		MsgTime: time.Now(),
	}

	if clientSendTime != "" {
		clientSendTimeStr := cast.ToInt64(clientSendTime)
		msg.ClientSendTime = time.Unix(0, clientSendTimeStr)
	}

	return msg
}
