package base

import (
	"github.com/summerKK/go-code-snippet-library/webcrawler/module/data"
)

type IDownloader interface {
	IModule
	Download(req *data.Request) (*data.Response, error)
}
