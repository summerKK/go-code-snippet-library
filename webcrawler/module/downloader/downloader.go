package downloader

import (
	"github.com/summerKK/go-code-snippet-library/webcrawler/logger"
	"github.com/summerKK/go-code-snippet-library/webcrawler/module"
	"github.com/summerKK/go-code-snippet-library/webcrawler/module/base"
	"net/http"
)

type Downloader struct {
	*module.Module
	// 下载用的client
	httpClient *http.Client
}

func New(mid base.MID, scoreCalc base.CalculateScore, httpClient *http.Client) (downloader *Downloader, err error) {

	m, err := module.NewModuleInternal(mid, scoreCalc)
	if err != nil {
		return nil, err
	}
	if httpClient == nil {
		return nil, genParameterError("nil http client")
	}

	return &Downloader{
		Module:     m,
		httpClient: httpClient,
	}, nil
}

func (d *Downloader) Download(req *base.Request) (resp *base.Response, err error) {
	d.Module.IncrHandlingNum()
	defer d.Module.DecrHandlingNum()
	d.Module.IncrCalledCount()
	if req == nil {
		err = genParameterError("nil request")
		return
	}
	httpReq := req.Req()
	if httpReq == nil {
		err = genParameterError("nil http request")
		return
	}
	d.Module.IncrAcceptedCount()
	logger.Logger.Infof("Do the request (URL:%s,depth:%d)", httpReq.URL, req.Depth())
	response, err := d.httpClient.Do(httpReq)
	if err != nil {
		return
	}
	d.Module.IncrCompletedCount()
	resp = base.NewResponse(response, req.Depth())

	return
}
