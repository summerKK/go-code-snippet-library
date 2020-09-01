package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/summerKK/go-code-snippet-library/chatroom/logic"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

var (
	userNum       int
	loginInterval time.Duration
	msgInterval   time.Duration
)

func init() {
	flag.IntVar(&userNum, "u", 100, "用户数量")
	flag.DurationVar(&loginInterval, "ll", 5e9, "登录间隔时间")
	flag.DurationVar(&msgInterval, "ml", time.Minute, "消息发送间隔时间")
}

func main() {
	flag.Parse()

	for i := 0; i < userNum; i++ {
		go UserConnect("user" + strconv.Itoa(i))
		time.Sleep(loginInterval)
	}

	select {}
}

func UserConnect(nickname string) {
	// 设置超时1分钟(测试一分钟)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, "ws://127.0.0.1:2020/ws?nickname="+nickname, nil)
	if err != nil {
		log.Println("Dial error:", err)
		return
	}

	log.Printf("user:%s登录", nickname)

	defer conn.Close(websocket.StatusInternalError, "内部错误!")

	go sendMessage(conn, nickname)

	for {
		var message logic.Message
		err := wsjson.Read(context.Background(), conn, &message)
		if err != nil {
			log.Println("receive msg error:", err)
			break
		}

		if message.ClientSendTime.IsZero() {
			continue
		}

		// 只打印延迟1s以上的消息
		if d := time.Now().Sub(message.ClientSendTime); d > 1*time.Second {
			fmt.Printf("接收到服务器响应(%d):%v\n", d.Milliseconds(), message)
		}
	}

	conn.Close(websocket.StatusNormalClosure, "")
}

func sendMessage(conn *websocket.Conn, nickname string) {
	ctx := context.Background()
	i := 1
	for {
		msg := map[string]string{
			"content":   fmt.Sprintf("来自%s的消息:%d", nickname, i),
			"send_time": strconv.FormatInt(time.Now().UnixNano(), 10),
		}
		err := wsjson.Write(ctx, conn, msg)
		if err != nil {
			log.Printf("send msg error:%v,nickname:%s,no:%d", err, nickname, i)
			break
		}
		i++

		time.Sleep(msgInterval)
	}
}
