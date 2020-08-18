package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, "ws://laravel-admin.work:2021/ws", nil)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(websocket.StatusInternalError, "内部错误")

	err = wsjson.Write(ctx, conn, "hello websocket service")
	if err != nil {
		log.Fatal(err)
	}

	var v interface{}
	err = wsjson.Read(ctx, conn, &v)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("接收到服务器响应:", v)

	conn.Close(websocket.StatusNormalClosure, "")

}
