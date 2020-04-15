package module

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

type Item map[string]interface{}

func (i *Item) Valid() bool {
	return i != nil
}
