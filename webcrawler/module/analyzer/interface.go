package analyzer

import (
	"github.com/summerKK/go-code-snippet-library/webcrawler/module"
	"github.com/summerKK/go-code-snippet-library/webcrawler/module/data"
	"net/http"
)

type ParseResponse func(httpResp *http.Response, respDepth uint32) ([]data.IData, error)

type IAnalyzer interface {
	module.IModule
	// 返回分析器使用的响应解析函数列表
	RespParsers() []ParseResponse
	Analyze(resp *data.Response) ([]data.IData, error)
}
