package errors

import (
	"bytes"
	"fmt"
)

type CrawlerError struct {
	errType    ErrorType
	errMsg     string
	fullErrMsg string
}

func (c *CrawlerError) Type() ErrorType {
	return c.errType
}

func (c *CrawlerError) Error() string {
	if c.fullErrMsg == "" {
		c.genFullErrMsg()
	}

	return c.fullErrMsg
}

func (c *CrawlerError) genFullErrMsg() {
	var buf bytes.Buffer
	buf.WriteString("crawler error: ")
	if c.errType != "" {
		buf.WriteString(string(c.errType))
		buf.WriteString(": ")
	}
	buf.WriteString(c.errMsg)
	c.fullErrMsg = fmt.Sprintf("%s", buf.String())
	return
}
