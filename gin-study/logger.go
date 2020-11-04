package gin

import (
	"log"
	"time"
)

func Logger() HandlerFunc {
	return func(c *Context) {
		t := time.Now()

		c.Next()

		log.Printf("[%d] %s in %v", c.Writer.Status(), c.Req.RequestURI, time.Since(t))
	}
}
