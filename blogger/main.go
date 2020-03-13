package main

import (
	"github.com/gin-gonic/gin"
	"summer/blogger/controller"
	"summer/blogger/dal/db"
)

var router *gin.Engine

func init() {
	router = gin.Default()
	dns := "root:root@tcp(127.0.0.1)/blogger?parseTime=true"
	err := db.Init(dns)
	if err != nil {
		panic(err)
	}
}

func main() {
	router.Static("/static/", "./static")
	router.LoadHTMLGlob("views/*")
}

func routing() {
	router.GET("/", controller.IndexHandler)
}
