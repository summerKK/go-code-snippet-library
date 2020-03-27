package main

import (
	"github.com/summerKK/go-code-snippet-library/spark/dal/db"
	"github.com/summerKK/go-code-snippet-library/spark/filter"
	"github.com/summerKK/go-code-snippet-library/spark/gin"
	idgen "github.com/summerKK/go-code-snippet-library/spark/id-gen"
	"github.com/summerKK/go-code-snippet-library/spark/logger"
	"github.com/summerKK/go-code-snippet-library/spark/middleware/account"
)

func init() {
	dns := "root:root@tcp(127.0.0.1)/spark?parseTime=true"
	// 初始化mysql
	err := db.Init(dns)
	if err != nil {
		panic(err)
	}
	// 初始化id生成器
	idgen.Init(0)
	// 初始化session
	err = account.InitSession()
	if err != nil {
		panic(err)
	}
	// 初始化敏感词过滤组件
	err = filter.Init("./data/filter.dat")
	if err != nil {
		panic(err)
	}
	// 初始化日志组件
	err = logger.Init()
	if err != nil {
		panic(err)
	}
	// 初始化gin框架
	err = gin.Init()
	if err != nil {
		panic(err)
	}
}

func main() {
}
