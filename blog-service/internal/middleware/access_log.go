package middleware

import (
	"bytes"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/summerKK/go-code-snippet-library/blog-service/global"
	"github.com/summerKK/go-code-snippet-library/blog-service/pkg/logger"
)

type AccessLogWrite struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w AccessLogWrite) Write(p []byte) (int, error) {
	if n, err := w.body.Write(p); err != nil {
		return n, err
	}

	return w.ResponseWriter.Write(p)
}

// 记录响应结果
func AccessLog() gin.HandlerFunc {
	return func(context *gin.Context) {
		bodyWrite := &AccessLogWrite{
			context.Writer,
			bytes.NewBufferString(""),
		}
		beginTime := time.Now().Unix()
		context.Next()
		endTime := time.Now().Unix()

		fields := logger.Fields{
			"request":  context.Request.PostForm.Encode(),
			"response": bodyWrite.body.String(),
		}

		global.Logger.WithFields(fields).Infof(
			context,
			"access log:method:%s,status_code:%d,begin_time:%d,end_time:%d",
			context.Request.Method,
			bodyWrite.Status(),
			beginTime,
			endTime,
		)
	}
}
