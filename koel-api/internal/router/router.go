package router

import (
	"github.com/gin-gonic/gin"
	"github.com/summerKK/go-code-snippet-library/koel-api/internal/middleware"
	"github.com/summerKK/go-code-snippet-library/koel-api/internal/router/api"
)

func NewRouter() *gin.Engine {
	r := gin.New()

	r.Use(middleware.Translations())

	r1 := r.Group("/api")
	r1.POST("/api/me", api.GetAuth)

	return r
}
