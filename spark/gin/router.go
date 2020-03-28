package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/summerKK/go-code-snippet-library/spark/controller/account"
	"github.com/summerKK/go-code-snippet-library/spark/controller/category"
	"github.com/summerKK/go-code-snippet-library/spark/controller/question"
	"github.com/summerKK/go-code-snippet-library/spark/middleware"
	middlewareAccount "github.com/summerKK/go-code-snippet-library/spark/middleware/account"
)

var engine *gin.Engine

func Init() (err error) {
	engine = gin.Default()
	router()
	return engine.Run("127.0.0.1:9080")
}

func router() {
	engine.Use(middleware.Cors())
	api := engine.Group("/api")
	api.POST("/user/register", account.Register)
	api.POST("/user/login", account.Login)

	group := api.Group("")
	group.Use(middlewareAccount.Auth())
	group.POST("/question/submit", question.Save)
	group.GET("/question/list", question.List)
	group.GET("/category/list", category.List)
}
