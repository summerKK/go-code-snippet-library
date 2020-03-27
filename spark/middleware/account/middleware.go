package account

import (
	"github.com/gin-gonic/gin"
	"github.com/summerKK/go-code-snippet-library/spark/util"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		ProcessRequest(c)
		if login, _ := IsLogin(c); !login {
			util.ResponseErr(c, &util.CodeInfo{Code: util.ErrUserNotLogin})
			c.Abort()
			return
		}
		c.Next()
	}
}
