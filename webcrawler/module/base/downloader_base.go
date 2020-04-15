package base

import (
	"github.com/summerKK/go-code-snippet-library/webcrawler/module"
)

type IDownloader interface {
	IModule
	Download(req *module.Request) (*module.Response, error)
}
