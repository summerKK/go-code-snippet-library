package analyzer

import (
	"github.com/summerKK/go-code-snippet-library/webcrawler/errors"
)

func genParameterError(errMsg string) error {
	return errors.NewCrawlerErrorBy(errors.ERROR_TYPE_ANALYZER,
		errors.NewIllegalParamsError(errMsg),
	)
}
