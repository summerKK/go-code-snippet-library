package analyzer

import (
	"fmt"
	"github.com/summerKK/go-code-snippet-library/webcrawler/module"
	"github.com/summerKK/go-code-snippet-library/webcrawler/module/base"
	"github.com/summerKK/go-code-snippet-library/webcrawler/module/data"
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

func (a *Analyzer) Analyze(resp *data.Response) ([]base.IData, error) {
	panic(nil)
}
