package base

import (
	"github.com/summerKK/go-code-snippet-library/webcrawler/module"
	"net/http"
)

type ParseResponse func(httpResp *http.Response, respDepth uint32) ([]IData, []error)

type IAnalyzer interface {
	IModule
	// 返回分析器使用的响应解析函数列表
	RespParsers() []ParseResponse
	Analyze(resp *module.Response) (datalist []IData, errlist []error)
}
