package main

import (
	"log"
	"net/http"

	"github.com/summerKK/go-code-snippet-library/chatroom/server"
)

var (
	addr = ":2020"
)

func main() {

	server.RegisterHandle()

	log.Fatal(http.ListenAndServe(addr, nil))
}
