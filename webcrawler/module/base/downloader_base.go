package base

type IDownloader interface {
	IModule
	Download(req *Request) (*Response, error)
}
