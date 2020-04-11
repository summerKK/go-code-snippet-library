package data

import "net/http"

type Request struct {
	req   *http.Request
	depth uint32
}

func (r *Request) Valid() bool {
	return r.req != nil && r.req.URL != nil
}

func (r *Request) Depth() uint32 {
	return r.depth
}

func (r *Request) Req() *http.Request {
	return r.req
}

func NewRequest(req *http.Request, depth uint32) *Request {
	return &Request{req: req, depth: depth}
}
