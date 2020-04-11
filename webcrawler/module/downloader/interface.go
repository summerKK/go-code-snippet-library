package downloader

import (
	"github.com/summerKK/go-code-snippet-library/webcrawler/module"
	"github.com/summerKK/go-code-snippet-library/webcrawler/module/data"
)

type IDownloader interface {
	module.IModule
	Download(req *data.Request) (*data.Response, error)
}
