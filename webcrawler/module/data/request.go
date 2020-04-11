package data

import "net/http"

type request struct {
	req   *http.Request
	depth uint32
}

func (r *request) Valid() bool {
	return r.req != nil && r.req.URL != nil
}

func (r *request) Depth() uint32 {
	return r.depth
}

func (r *request) Req() *http.Request {
	return r.req
}

func NewRequest(req *http.Request, depth uint32) *request {
	return &request{req: req, depth: depth}
}
