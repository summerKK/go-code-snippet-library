package errors

import (
	"fmt"
	"testing"
)

func TestCrawlerError(t *testing.T) {
	crawlerError := NewCrawlerError(ERROR_TYPE_ANALYZER, "hello,world")
	fmt.Println(crawlerError)
}
