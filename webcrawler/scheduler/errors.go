package scheduler

import (
	"github.com/summerKK/go-code-snippet-library/webcrawler/errors"
	"github.com/summerKK/go-code-snippet-library/webcrawler/logger"
	"github.com/summerKK/go-code-snippet-library/webcrawler/module"
	"github.com/summerKK/go-code-snippet-library/webcrawler/module/base"
	"github.com/summerKK/go-code-snippet-library/webcrawler/toolkit/buffer"
)

// genError 用于生成爬虫错误值。
func genError(errMsg string) error {
	return errors.NewCrawlerError(errors.ERROR_TYPE_SCHEDULER,
		errMsg)
}

// genErrorByError 用于基于给定的错误值生成爬虫错误值。
func genErrorByError(err error) error {
	return errors.NewCrawlerError(errors.ERROR_TYPE_SCHEDULER,
		err.Error())
}

// sendError 用于向错误缓冲池发送错误值。
func sendError(err error, mid base.MID, errorBufferPool buffer.IPool) bool {
	if err == nil || errorBufferPool == nil || errorBufferPool.Closed() {
		return false
	}
	var crawlerError errors.ICrawlerError
	var ok bool
	crawlerError, ok = err.(errors.ICrawlerError)
	if !ok {
		var moduleType base.MType
		var errorType errors.ErrorType
		ok, moduleType = module.GetType(mid)
		if !ok {
			errorType = errors.ERROR_TYPE_SCHEDULER
		} else {
			switch moduleType {
			case base.TYPE_DOWNLOADER:
				errorType = errors.ERROR_TYPE_DOWNLOADER
			case base.TYPE_ANALYZER:
				errorType = errors.ERROR_TYPE_ANALYZER
			case base.TYPE_PIPELINE:
				errorType = errors.ERROR_TYPE_PIPELINE
			}
		}
		crawlerError = errors.NewCrawlerError(errorType, err.Error())
	}
	if errorBufferPool.Closed() {
		return false
	}
	go func(crawlerError errors.ICrawlerError) {
		if err := errorBufferPool.Put(crawlerError); err != nil {
			logger.Logger.Warnln("The error buffer pool was closed. Ignore error sending.")
		}
	}(crawlerError)

	return true
}
