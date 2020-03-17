package main

import (
	"github.com/DeanThompson/ginpprof"
	"github.com/gin-gonic/gin"
	"github.com/summerKK/go-code-snippet-library/blogger/controller"
	"github.com/summerKK/go-code-snippet-library/blogger/dal/db"
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

	ginpprof.Wrap(router)
}
