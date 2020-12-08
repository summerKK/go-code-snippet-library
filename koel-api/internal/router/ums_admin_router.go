package router

import (
	"github.com/gin-gonic/gin"
	"github.com/summerKK/go-code-snippet-library/koel-api/internal/router/api"
)

func umsAdminRoute(r *gin.RouterGroup) {
	group := r.Group("/admin")

	group.POST("/login", api.GetAuth)
}
