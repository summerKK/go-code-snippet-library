package server

import (
	"log"
	"net/http"

	"github.com/summerKK/go-code-snippet-library/chatroom/logic"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func websocketHandleFunc(w http.ResponseWriter, r *http.Request) {
	// Accept 从客户端接受 WebSocket 握手，并将连接升级到 WebSocket。
	// 如果 Origin 域与主机不同，Accept 将拒绝握手，除非设置了 InsecureSkipVerify 选项（通过第三个参数 AcceptOptions 设置）。
	// 换句话说，默认情况下，它不允许跨源请求。如果发生错误，Accept 将始终写入适当的响应
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{InsecureSkipVerify: true})
	if err != nil {
		log.Printf("websocket.Accept error:%v", err)
		return
	}

	token := r.FormValue("token")
	nickname := r.FormValue("nickname")
	if l := len(nickname); l < 2 || l > 20 {
		log.Println("illegal nickname")
		wsjson.Write(r.Context(), conn, logic.NewErrorMsg("非法昵称，昵称长度：2-20"))
		conn.Close(websocket.StatusUnsupportedData, "illegal nickname")
		return
	}

	if !logic.Broadcaster.CanEnterRoom(nickname) {
		log.Println("昵称已存在:", nickname)
		wsjson.Write(r.Context(), conn, logic.NewErrorMsg("该昵称已经存在!"))
		conn.Close(websocket.StatusUnsupportedData, "该昵称已经存在")
		return
	}

	user := logic.NewUser(nickname, token, r.RemoteAddr, conn)

	// 开启给用户发送消息的goroutine
	go user.SendMessage(r.Context())

	// 给新用户发送欢迎语
	user.MessageChannel <- logic.NewWelcomeMsg(user)

	// 提醒所有用户新用户到来
	msg := logic.NewEnterMsg(user)
	logic.Broadcaster.Broadcast(msg)

	// 把用户加入到用户列表
	logic.Broadcaster.UserEntering(user)
	log.Printf("用户:%s加入聊天室\n", user.Nickname)

	// 给用户发送历史消息
	logic.OfflineProcessor.Send(user)

	// 用户接收消息
	err = user.ReceiveMessage(r.Context())

	// 用户断开连接(用户离开)
	logic.Broadcaster.UserLeaving(user)
	log.Printf("用户:%s离开聊天室\n", user.Nickname)
	msg = logic.NewLeavingMsg(user)
	logic.Broadcaster.Broadcast(msg)

	if err == nil {
		conn.Close(websocket.StatusNormalClosure, "")
	} else {
		log.Printf("read from client error:%v", err)
		conn.Close(websocket.StatusInternalError, "Read from client error")
	}
}
