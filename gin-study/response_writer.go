package gin

import "net/http"

type ResponseWriterInterface interface {
	http.ResponseWriter
	Status() int
	Written() bool

	SetStatus(int)
	Reset(w http.ResponseWriter)
}

type ResponseWriter struct {
	http.ResponseWriter
	status  int
	written bool
}

func NewResponseWriter(w http.ResponseWriter, status int, written bool) *ResponseWriter {
	return &ResponseWriter{
		w,
		status,
		written,
	}
}

func (r *ResponseWriter) Status() int {
	return r.status
}

func (r *ResponseWriter) Written() bool {
	return r.written
}

func (r *ResponseWriter) SetStatus(i int) {
	r.status = i
}

func (r *ResponseWriter) Reset(w http.ResponseWriter) {
	r.status = 0
	r.written = false
	r.ResponseWriter = w
}

func (r *ResponseWriter) WriteHeader(s int) {
	r.ResponseWriter.WriteHeader(s)
	r.status = s
}

func (r *ResponseWriter) Write(b []byte) (int, error) {
	return r.ResponseWriter.Write(b)
}
