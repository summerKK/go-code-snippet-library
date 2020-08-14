package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"strconv"
	"sync/atomic"
	"time"
)

type User struct {
	ID      int
	Addr    string
	EnterAt time.Time
	// 用户接受消息的channel
	MessageChannel chan string
}

func (u *User) String() string {
	return fmt.Sprintf("%d:%s", u.ID, u.Addr)
}

var (
	// 生成用户Id
	seqNo int64 = 0
	// 全局消息共享channel
	messageChannel chan []string = make(chan []string)
	// 进入chatroom的用户(进入中)
	enteringChannel chan *User = make(chan *User)
	// 正在离开chatroom的用户
	leavingChannel chan *User = make(chan *User)
)

func main() {

	go func() {
		http.ListenAndServe(":2021", nil)
	}()

	listener, err := net.Listen("tcp", ":2020")
	if err != nil {
		panic(err)
	}

	go broadcaster()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("listener.Accept error:%v", err)
			continue
		}

		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	user := &User{
		ID:             GenUserId(),
		Addr:           conn.RemoteAddr().String(),
		EnterAt:        time.Now(),
		MessageChannel: make(chan string, 8),
	}

	// 监听消息
	go sendMessage(conn, user.MessageChannel)

	// 给登录用户发送消息
	user.MessageChannel <- "welcome," + user.String()
	// 给其他在线用户发送消息
	messageChannel <- []string{"user:" + strconv.Itoa(user.ID) + " has enter", strconv.Itoa(user.ID)}

	enteringChannel <- user

	input := bufio.NewScanner(conn)
	for input.Scan() {
		messageChannel <- []string{strconv.Itoa(user.ID) + ":" + input.Text(), strconv.Itoa(user.ID)}
	}

	if err := input.Err(); err != nil {
		log.Printf("读取错误:%v", err)
	}

	// 用户离开,断开连接
	leavingChannel <- user
	messageChannel <- []string{"user:" + strconv.Itoa(user.ID) + " has left", strconv.Itoa(user.ID)}

}

// 给单个用户发送消息
func sendMessage(conn net.Conn, channel <-chan string) {
	for s := range channel {
		_, _ = fmt.Fprintln(conn, s)
	}
}

func GenUserId() int {
	atomic.AddInt64(&seqNo, 1)

	return int(atomic.LoadInt64(&seqNo))
}

func broadcaster() {
	// 在线用户
	users := make(map[*User]struct{}, 256)
	for {
		select {
		case user := <-enteringChannel:
			users[user] = struct{}{}
		case user := <-leavingChannel:
			delete(users, user)
			// 记住关闭channel,避免造成内存泄漏
			close(user.MessageChannel)
			// 必须使用goroutine,要不然select可能会阻塞.或者直接在 leavingChannel <- user 的时候传消息
			go func() {
				messageChannel <- []string{fmt.Sprintf("%d user离开chatroom", user.ID), strconv.Itoa(user.ID)}
			}()
		case s := <-messageChannel:
			// 给所有在线用户发送消息
			for user := range users {
				//  发送消息用户不接收消息
				if strconv.Itoa(user.ID) == s[1] {
					continue
				}

				user.MessageChannel <- s[0]
			}
		}
	}
}
