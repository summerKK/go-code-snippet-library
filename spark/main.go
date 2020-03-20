package main

import (
	"github.com/summerKK/go-code-snippet-library/spark/dal/db"
	"github.com/summerKK/go-code-snippet-library/spark/gin"
)

func init() {
	dns := "root:root@tcp(127.0.0.1)/spark?parseTime=true"
	// 初始化mysql
	err := db.Init(dns)
	if err != nil {
		panic(err)
	}
	// 初始化gin框架
	gin.Init()
}

func main() {
}
