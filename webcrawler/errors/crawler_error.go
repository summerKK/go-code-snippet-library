package errors

import (
	"bytes"
	"fmt"
	"runtime"
)

type CrawlerError struct {
	errType    ErrorType
	errMsg     string
	fullErrMsg string
}

func NewCrawlerError(errType ErrorType, errMsg string) *CrawlerError {
	return &CrawlerError{errType: errType, errMsg: errMsg}
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
	_, file, line, ok := runtime.Caller(1)
	if ok {
		buf.WriteByte('\n')
		buf.WriteString(fmt.Sprintf("     file:%s ", file))
		buf.WriteString(fmt.Sprintf(" line:%d ", line))
	}

	c.fullErrMsg = fmt.Sprintf("%s", buf.String())
	return
}
