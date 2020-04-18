package pipeline

import "github.com/summerKK/go-code-snippet-library/webcrawler/errors"

func genParameterError(msg string) error {
	return errors.NewCrawlerErrorBy(errors.ERROR_TYPE_PIPELINE, errors.NewIllegalParamsError(msg))
}

func genError(msg string) error {
	return errors.NewCrawlerError(errors.ERROR_TYPE_PIPELINE, msg)
}
