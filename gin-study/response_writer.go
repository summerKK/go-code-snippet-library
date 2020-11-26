package gin

import (
	"bufio"
	"errors"
	"log"
	"net"
	"net/http"
)

type ResponseWriterInterface interface {
	Status() int
	Written() bool
	WriteHeaderNow()
	Reset(w http.ResponseWriter)

	http.ResponseWriter
	http.Hijacker
	http.Flusher
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
	if s > 0 {
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

// 劫持
func (r *ResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := r.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("the ResponseWriter doesn't support the Hijacker interface")
	}

	return hijacker.Hijack()
}

// implements the http.Flush interface
func (r *ResponseWriter) Flush() {
	if flusher, ok := r.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}
