package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/summerKK/go-code-snippet-library/blog-service/pkg/app"
	"github.com/summerKK/go-code-snippet-library/blog-service/pkg/errcode"
	"github.com/summerKK/go-code-snippet-library/blog-service/pkg/limiter"
)

func RateLimiter(m limiter.Contract) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := m.Key(c)
		if bucket, ok := m.GetBucket(key); ok {
			count := bucket.TakeAvailable(1)
			if count == 0 {
				response := app.NewResponse(c)
				response.ToErrorResponse(errcode.TooManyRequests)
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
