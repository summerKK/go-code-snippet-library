package downloader

import "github.com/summerKK/go-code-snippet-library/webcrawler/errors"

// genParameterError 用于生成爬虫参数错误值。
func genParameterError(errMsg string) error {
	return errors.NewCrawlerErrorBy(errors.ERROR_TYPE_DOWNLOADER,
		errors.NewIllegalParamsError(errMsg))
}
