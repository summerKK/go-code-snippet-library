package data

import "net/http"

type Response struct {
	resp  *http.Response
	depth uint32
}

func (r *Response) Valid() bool {
	return r.resp != nil && r.resp.Body != nil
}

func (r *Response) Depth() uint32 {
	return r.depth
}

func (r *Response) Resp() *http.Response {
	return r.resp
}

func NewResponse(resp *http.Response, depth uint32) *Response {
	return &Response{resp: resp, depth: depth}
}
