package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/summerKK/go-code-snippet-library/blog-service/global"
	"github.com/summerKK/go-code-snippet-library/blog-service/pkg/app"
	"github.com/summerKK/go-code-snippet-library/blog-service/pkg/email"
	"github.com/summerKK/go-code-snippet-library/blog-service/pkg/errcode"
)

func Recovery() gin.HandlerFunc {
	defaultMailer := email.NewEmail(&email.SMTPInfo{
		Host:     global.EmailSetting.Host,
		Port:     global.EmailSetting.Port,
		IsSSL:    global.EmailSetting.IsSSL,
		UserName: global.EmailSetting.UserName,
		Password: global.EmailSetting.Password,
		From:     global.EmailSetting.From,
	})

	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				global.Logger.WithCallerFrames().Errorf(c, "panic recover err:%v", err)

				if defaultMailer.UserName != "" && defaultMailer.Password != "" {
					err := defaultMailer.SendMail(
						global.EmailSetting.To,
						fmt.Sprintf("异常抛出,发生时间:%d", time.Now().Unix()),
						fmt.Sprintf("错误信息: %v", err),
					)
					if err != nil {
						global.Logger.Panicf(c, "mail.SendMail err:%v", err)
					}
				}

				app.NewResponse(c).ToErrorResponse(errcode.ServerError)
				c.Abort()
			}
		}()

		c.Next()
	}
}
