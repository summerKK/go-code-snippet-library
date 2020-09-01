package main

import (
	"log"
	"net/http"

	_ "net/http/pprof"

	"github.com/spf13/viper"
	_ "github.com/summerKK/go-code-snippet-library/chatroom/global"
	"github.com/summerKK/go-code-snippet-library/chatroom/server"
)

var (
	addr = viper.GetString("service-port")
)

func main() {

	server.RegisterHandle()

	log.Fatal(http.ListenAndServe(addr, nil))
}
