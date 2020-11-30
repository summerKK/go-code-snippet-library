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

		clientIp := c.ClientIp()
		method := c.Request.Method
		methodColor := colorForMethod(method)
		statusCode := c.Writer.Status()
		statusColor := colorForStatus(statusCode)

		logger.Printf("[GIN] %v |%s %3d %s| %12v | %s |%s %-7s %s|%s\n%s",
			time.Now().Format("2006/01/02 - 15:04:05"),
			statusColor, statusCode, reset,
			time.Since(t),
			clientIp,
			methodColor,
			method,
			reset,
			c.Request.URL.Path,
			c.Errors.String(),
		)
	}
}

func colorForStatus(code int) string {
	switch {
	case code >= 200 && code <= 299:
		return green
	case code >= 300 && code <= 399:
		return white
	case code >= 400 && code <= 499:
		return yellow
	default:
		return red
	}
}

func colorForMethod(method string) string {
	switch {
	case method == "GET":
		return blue
	case method == "POST":
		return cyan
	case method == "PUT":
		return yellow
	case method == "DELETE":
		return red
	case method == "PATCH":
		return green
	case method == "HEAD":
		return magenta
	case method == "OPTIONS":
		return white
	default:
		return reset
	}
}
