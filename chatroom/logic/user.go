package logic

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
	"sync/atomic"
	"time"

	"github.com/spf13/cast"
	"github.com/spf13/viper"
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

	// 标识当前用户是否已经存在
	isNew bool
}

var globalUid uint64 = 0

// 系统用户
var System = &User{}

// 解析@
var regAtUsers = regexp.MustCompile(`@[^\s@]{2,20}`)

func NewUser(nickname, token, addr string, conn *websocket.Conn) *User {
	user := &User{
		Nickname:       nickname,
		EnterAt:        time.Now(),
		Addr:           addr,
		MessageChannel: make(chan *Message, 32),
		Token:          token,
		conn:           conn,
	}

	// 已经登录过的用户
	if user.Token != "" {
		uid, err := parseTokenAndValidate(token, nickname)
		if err == nil {
			user.Uid = int(uid)
		}
	}

	if user.Uid == 0 {
		user.Uid = int(atomic.AddUint64(&globalUid, 1))
		user.Token = genToken(user.Uid, user.Nickname)
		user.isNew = true
	}

	return user
}

func parseTokenAndValidate(token string, nickname string) (uint64, error) {
	pos := strings.LastIndex(token, "uid")
	messageMac0, err := base64.StdEncoding.DecodeString(token[:pos])
	if err != nil {
		return 0, err
	}
	// pos + 3  去掉 `uid`  字符串
	uid := cast.ToUint64(token[pos+3:])

	secret := viper.GetString("token-secret")
	message := fmt.Sprintf("%s%s%d", nickname, secret, uid)
	messageMac1 := macSha256([]byte(message), []byte(secret))

	// 验证token是否正确
	ok := hmac.Equal(messageMac0, messageMac1)
	if ok {
		return uid, nil
	}

	return 0, errors.New("token is illegal")
}

func genToken(uid int, nickname string) string {
	secret := viper.GetString("token-secret")
	message := fmt.Sprintf("%s%s%d", nickname, secret, uid)
	messageMac := macSha256([]byte(message), []byte(secret))

	return fmt.Sprintf("%suid%d", base64.StdEncoding.EncodeToString(messageMac), uid)
}

func macSha256(message, secret []byte) []byte {
	mac := hmac.New(sha256.New, secret)
	mac.Write(message)

	return mac.Sum(nil)
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

		// 过滤敏感词
		receiveMsg["content"] = FilterSensitiveWord(receiveMsg["content"])

		// 内容发送到聊天室
		sendMsg := NewMsg(u, receiveMsg["content"], receiveMsg["client_time"])

		// 解析content,看看@了谁
		sendMsg.Ats = regAtUsers.FindAllString(sendMsg.Content, -1)

		Broadcaster.Broadcast(sendMsg)
	}
}
