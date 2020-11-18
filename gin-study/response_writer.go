package gin

import (
	"log"
	"net/http"
)

type ResponseWriterInterface interface {
	http.ResponseWriter
	Status() int
	Written() bool

	WriteHeaderNow()
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

func (r *ResponseWriter) Reset(w http.ResponseWriter) {
	r.status = http.StatusOK
	r.written = false
	r.ResponseWriter = w
}

func (r *ResponseWriter) WriteHeader(s int) {
	if s != 0 {
		r.status = s
		if r.written {
			log.Println("[GIN] WARNING. Headers were already written!")
		}
	}
}

func (r *ResponseWriter) WriteHeaderNow() {
	if !r.written {
		r.written = true
		r.ResponseWriter.WriteHeader(r.status)
	}
}

func (r *ResponseWriter) Write(b []byte) (int, error) {
	r.WriteHeaderNow()

	return r.ResponseWriter.Write(b)
}
