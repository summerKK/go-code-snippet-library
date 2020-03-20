package main

import (
	"github.com/DeanThompson/ginpprof"
	"github.com/gin-gonic/gin"
	"github.com/summerKK/go-code-snippet-library/spark/dal/db"
)

var router *gin.Engine

func init() {
	router = gin.Default()
	dns := "root:root@tcp(127.0.0.1)/spark?parseTime=true"
	err := db.Init(dns)
	if err != nil {
		panic(err)
	}
}

func main() {
	routing()
	router.Run()
}

func routing() {

	ginpprof.Wrap(router)
}
