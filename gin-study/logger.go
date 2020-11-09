package gin

import (
	"log"
	"time"
)

func ErrorLogger() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if len(c.Errors) > 0 {
				log.Println(c.Errors)
				c.JSON(-1, c.Errors)
			}
		}()

		c.Next()
	}
}

func Logger() HandlerFunc {
	return func(c *Context) {
		t := time.Now()

		c.Next()

		log.Printf("[%d] %s in %v", c.Writer.Status(), c.Req.RequestURI, time.Since(t))
	}
}
