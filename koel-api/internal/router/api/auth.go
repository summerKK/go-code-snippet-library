package api

import (
	"github.com/gin-gonic/gin"
	"github.com/summerKK/go-code-snippet-library/koel-api/internal/service"
	"github.com/summerKK/go-code-snippet-library/koel-api/pkg/app"
	"github.com/summerKK/go-code-snippet-library/koel-api/pkg/errcode"
)

func GetAuth(c *gin.Context) {
	params := &service.UserRequest{}
	response := app.NewResponse(c)
	ok, errors := app.BindAndValid(c, params)
	if !ok {
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errors.Errors()...))
		return
	}

	svc := service.NewService(c.Request.Context())
	err := svc.CheckAuth(params)
	if err != nil {
		response.ToErrorResponse(errcode.UnauthorizedAuthNotExist)
		return
	}

	token, err := app.GenerateToken(params.Email)
	if err != nil {
		response.ToErrorResponse(errcode.UnauthorizedTokenGenerate)
		return
	}

	response.ToResponse(gin.H{
		"token": token,
	})
}
