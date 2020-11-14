package gin

import (
	"io"
	"log"
	"os"
	"time"
)

func ErrorLogger() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if len(c.Errors) > 0 {
				c.JSON(-1, c.Errors)
			}
		}()

		c.Next()
	}
}

var (
	green  = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	white  = string([]byte{27, 91, 57, 48, 59, 52, 55, 109})
	yellow = string([]byte{27, 91, 57, 55, 59, 52, 51, 109})
	red    = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	reset  = string([]byte{27, 91, 48, 109})
)

func Logger(writer io.Writer) HandlerFunc {
	if writer == nil {
		writer = os.Stdout
	}
	logger := log.New(writer, "", 0)

	return func(c *Context) {
		t := time.Now()

		c.Next()

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

		logger.Printf("[GIN] %v |%s %3d %s| %12v | %3.1f%% | %3s %s\n",
			time.Now().Format("2006/01/02 - 15:04:05"),
			color, c.Writer.Status(), reset,
			time.Since(t),
			c.Engine.CacheStress()*100,
			c.Req.Method, c.Req.URL.Path,
		)

		if len(c.Errors) > 0 {
			log.Println(c.Errors)
		}
	}
}
