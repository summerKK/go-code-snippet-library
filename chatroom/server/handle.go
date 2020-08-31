package server

import (
	"net/http"

	"github.com/summerKK/go-code-snippet-library/chatroom/logic"
)

func RegisterHandle() {

	// 广播消息处理
	go logic.Broadcaster.Start()

	http.HandleFunc("/", homeHandleFunc)
	http.HandleFunc("/ws", websocketHandleFunc)
}
