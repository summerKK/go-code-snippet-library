package gin

import (
	"io"
	"log"
	"os"
	"time"
)

var (
	green   = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	white   = string([]byte{27, 91, 57, 48, 59, 52, 55, 109})
	yellow  = string([]byte{27, 91, 57, 55, 59, 52, 51, 109})
	red     = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	blue    = string([]byte{27, 91, 57, 55, 59, 52, 52, 109})
	magenta = string([]byte{27, 91, 57, 55, 59, 52, 53, 109})
	cyan    = string([]byte{27, 91, 57, 55, 59, 52, 54, 109})
	reset   = string([]byte{27, 91, 48, 109})
)

// 返回所有错误
func ErrorLogger() HandlerFunc {
	return ErrorLoggerT(ErrorTypeAll)
}

//  返回特定的错误
func ErrorLoggerT(typ uint32) HandlerFunc {
	return func(c *Context) {

		c.Next()

		errs := c.Errors.ByType(typ)
		if len(errs) > 0 {
			c.JSON(-1, c.Errors)
		}
	}
}

func Logger(writer io.Writer) HandlerFunc {
	if writer == nil {
		writer = os.Stdout
	}
	logger := log.New(writer, "", 0)

	return func(c *Context) {
		t := time.Now()

		c.Next()

		requester := c.Request.Header.Get("X-Real-IP")
		if requester == "" {
			requester = c.Request.Header.Get("X-Forwarded-For")
		}

		// 如果还是为空直接取request的ip
		if requester == "" {
			requester = c.Request.RemoteAddr
		}

		var color string
		code := c.Writer.Status()
		switch {
		case code >= 200 && code <= 299:
			color = green
		case code >= 300 && code <= 399:
			color = white
		case code >= 400 && code <= 499:
			color = yellow
		default:
			color = red
		}

		var methodColor string
		method := c.Request.Method
		switch {
		case method == "GET":
			methodColor = blue
		case method == "POST":
			methodColor = cyan
		case method == "PUT":
			methodColor = yellow
		case method == "DELETE":
			methodColor = red
		case method == "PATCH":
			methodColor = green
		case method == "HEAD":
			methodColor = magenta
		case method == "OPTIONS":
			methodColor = white
		}

		logger.Printf("[GIN] %v |%s %3d %s| %12v | %s |%s %-7s %s|%s\n%s",
			time.Now().Format("2006/01/02 - 15:04:05"),
			color, c.Writer.Status(), reset,
			time.Since(t),
			requester,
			methodColor,
			c.Request.Method,
			reset,
			c.Request.URL.Path,
			c.Errors.String(),
		)
	}
}
