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
	routing()

	router.Run()
}

func routing() {
	router.GET("/", controller.IndexHandler)
	router.GET("/article/new", controller.ArticleCreate)
	router.POST("/article/submit", controller.ArticleSubmit)
	router.GET("/article/detail/", controller.ArticleInfo)
	router.POST("/article/comment/submit/", controller.CommentSubmit)
}
