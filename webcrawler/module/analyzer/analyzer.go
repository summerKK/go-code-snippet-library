package analyzer

import (
	"fmt"
	"github.com/summerKK/go-code-snippet-library/webcrawler/logger"
	"github.com/summerKK/go-code-snippet-library/webcrawler/module"
	"github.com/summerKK/go-code-snippet-library/webcrawler/module/base"
	"github.com/summerKK/go-code-snippet-library/webcrawler/toolkit/reader"
)

type Analyzer struct {
	*module.Module
	respParsers []base.ParseResponse
}

func New(mid base.MID, scoreCalc base.CalculateScore, respParsers []base.ParseResponse) (analyzer *Analyzer, err error) {
	m, err := module.NewModuleInternal(mid, scoreCalc)
	if err != nil {
		return
	}
	if respParsers == nil {
		err = genParameterError("nil respParsers")
		return
	}
	if len(respParsers) == 0 {
		err = genParameterError("empty respParsers")
		return
	}

	// 主要是为了防止在分析器创建后外界再对解析器列表进行更改.所以赋值给新的变量
	var innerParsers []base.ParseResponse
	for i, parser := range respParsers {
		if parser == nil {
			err = genParameterError(fmt.Sprintf("nil response parse [%d]", i))
			return
		}
		innerParsers = append(innerParsers, parser)
	}

	return &Analyzer{
		respParsers: innerParsers,
		Module:      m,
	}, nil
}

func (a *Analyzer) RespParsers() []base.ParseResponse {
	return a.respParsers
}

func (a *Analyzer) Analyze(resp *module.Response) (datalist []base.IData, errlist []error) {
	a.Module.IncrHandlingNum()
	defer a.Module.DecrHandlingNum()
	a.Module.IncrCalledCount()
	if resp == nil {
		errlist = append(errlist, genParameterError("nil resp"))
		return
	}
	httpResp := resp.Resp()
	if httpResp == nil {
		errlist = append(errlist, genParameterError("nil http response"))
		return
	}
	httpReq := httpResp.Request
	if httpReq == nil {
		errlist = append(errlist, genParameterError("nil http request"))
		return
	}
	reqUrl := httpReq.URL
	if reqUrl == nil {
		errlist = append(errlist, genParameterError("nil http request url"))
		return
	}
	a.Module.IncrAcceptedCount()
	respDepth := resp.Depth()
	logger.Logger.Infof("Parse the response (URL:%s,depth:%d)", reqUrl, respDepth)

	// 这里要把原生的httpResponse.Body保存起来
	originalRespBody := httpResp.Body
	if originalRespBody != nil {
		defer originalRespBody.Close()
	}
	// 这个reader不会关闭`reader`,意思就是可以重复被读取
	multipleReader, err := reader.NewReader(httpResp.Body)
	if err != nil {
		errlist = append(errlist, genError(err.Error()))
		return
	}

	for _, parser := range a.respParsers {
		// 创建一个NoCloser().
		httpResp.Body = multipleReader.Reader()
		pDatalist, pErrlist := parser(httpResp, respDepth)
		if pDatalist != nil {
			for _, pdata := range pDatalist {
				if pdata == nil {
					continue
				}
				datalist = appendDataList(datalist, pdata, respDepth)
			}
		}
		if pErrlist != nil {
			for _, pErr := range pErrlist {
				if pErr == nil {
					continue
				}
				errlist = append(errlist, pErr)
			}
		}
		if len(pErrlist) == 0 {
			a.Module.IncrCompletedCount()
		}
	}

	return
}

// appendDataList 用于添加请求值或条目值到列表。
func appendDataList(datalist []base.IData, data base.IData, respDepth uint32) []base.IData {
	if data == nil {
		return datalist
	}
	req, ok := data.(*module.Request)
	if !ok {
		datalist = append(datalist, data)
		return datalist
	}
	newDepth := respDepth + 1
	if newDepth != req.Depth() {
		req = module.NewRequest(req.Req(), newDepth)
	}

	return append(datalist, req)
}
