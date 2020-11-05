package main

import (
	"context"
	"log"

	"github.com/summerKK/go-code-snippet-library/gin-study"
)

func main() {
	engine := gin.Default()

	engine.POST("/api/user", func(c *gin.Context) {
		c.Writer.Write([]byte("hello,world"))
		log.Println("api/user")
	})

	engine.Run(context.Background(), ":8080")
}
