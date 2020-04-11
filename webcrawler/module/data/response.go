package data

import "net/http"

type response struct {
	resp  *http.Response
	depth uint32
}

func (r *response) Valid() bool {
	return r.resp != nil && r.resp.Body != nil
}

func (r *response) Depth() uint32 {
	return r.depth
}

func (r *response) Resp() *http.Response {
	return r.resp
}

func NewResponse(resp *http.Response, depth uint32) *response {
	return &response{resp: resp, depth: depth}
}
