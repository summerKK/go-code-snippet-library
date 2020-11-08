package main

import (
	"context"

	"github.com/summerKK/go-code-snippet-library/gin-study"
)

func main() {
	engine := gin.Default()

	engine.POST("/api/user", func(c *gin.Context) {
		c.Abort(401)
	})

	engine.Run(context.Background(), ":8080")
}
