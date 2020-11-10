package gin

import "net/http"

type ResponseWriter interface {
	http.ResponseWriter
	Status() int
	Written() bool

	SetStatus(int)
	Reset(w http.ResponseWriter)
}

type responseWriter struct {
	http.ResponseWriter
	status  int
	written bool
}

func (r *responseWriter) Status() int {
	return r.status
}

func (r *responseWriter) Written() bool {
	return r.written
}

func (r *responseWriter) SetStatus(i int) {
	r.status = i
}

func (r *responseWriter) Reset(w http.ResponseWriter) {
	r.status = 0
	r.written = false
	r.ResponseWriter = w
}

func (r *responseWriter) WriteHeader(s int) {
	r.ResponseWriter.WriteHeader(s)
	r.status = s
}

func (r *responseWriter) Write(b []byte) (int, error) {
	return r.ResponseWriter.Write(b)
}
