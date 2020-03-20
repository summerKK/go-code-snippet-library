package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/summerKK/go-code-snippet-library/spark/controller/account"
	"github.com/summerKK/go-code-snippet-library/spark/middleware"
)

var engine *gin.Engine

func Init() {
	engine = gin.Default()
	router()
	engine.Run(":9090")
}

func router() {
	engine.Use(middleware.Cors())
	api := engine.Group("/api")
	api.POST("/user/register", account.Register)
}
