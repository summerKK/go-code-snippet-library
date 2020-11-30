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
	Size() int
	Reset(w http.ResponseWriter)

	http.ResponseWriter
	http.Hijacker
	http.Flusher
}

const (
	NoWritten = -1
)

type ResponseWriter struct {
	http.ResponseWriter
	status int
	size   int
}

func NewResponseWriter(w http.ResponseWriter, status int) *ResponseWriter {
	return &ResponseWriter{
		w,
		status,
		NoWritten,
	}
}

func (r *ResponseWriter) Status() int {
	return r.status
}

func (r *ResponseWriter) Size() int {
	return r.size
}

func (r *ResponseWriter) Written() bool {
	return r.size != NoWritten
}

func (r *ResponseWriter) Reset(w http.ResponseWriter) {
	r.status = http.StatusOK
	r.size = NoWritten
	r.ResponseWriter = w
}

func (r *ResponseWriter) WriteHeader(s int) {
	if s > 0 {
		r.status = s
		if r.Written() {
			log.Println("[GIN] WARNING. Headers were already written!")
		}
	}
}

func (r *ResponseWriter) WriteHeaderNow() {
	if !r.Written() {
		r.size = 0
		r.ResponseWriter.WriteHeader(r.status)
	}
}

func (r *ResponseWriter) Write(b []byte) (int, error) {
	r.WriteHeaderNow()

	n, err := r.ResponseWriter.Write(b)
	r.size += n

	return n, err
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
