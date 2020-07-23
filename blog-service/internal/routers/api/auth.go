package api

import (
	"github.com/gin-gonic/gin"
	"github.com/summerKK/go-code-snippet-library/blog-service/global"
	"github.com/summerKK/go-code-snippet-library/blog-service/internal/service"
	"github.com/summerKK/go-code-snippet-library/blog-service/pkg/app"
	"github.com/summerKK/go-code-snippet-library/blog-service/pkg/errcode"
)

func GetAuth(c *gin.Context) {
	param := &service.AuthRequest{}
	response := app.NewResponse(c)
	ok, errors := app.BindAndValid(c, param)
	if ok {
		global.Logger.Errorf(c, "app.BindAndValid error;%v", errors)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errors.Errors()...))
		return
	}

	svc := service.New(c.Request.Context())
	err := svc.CheckAuth(param)
	if err != nil {
		global.Logger.Errorf(c, "svc.checkAuth error:%v", err)
		response.ToErrorResponse(errcode.UnauthorizedAuthNotExist)
		return
	}

	token, err := app.GenerateToken(param.AppKey, param.AppSecret)
	if err != nil {
		global.Logger.Errorf(c, "app.GenerateToken error:%v", err)
		response.ToErrorResponse(errcode.UnauthorizedTokenGenerate)
		return
	}

	response.ToResponse(gin.H{
		"token": token,
	})

	return
}
