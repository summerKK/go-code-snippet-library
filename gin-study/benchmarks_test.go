package gin_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/summerKK/go-code-snippet-library/gin-study"
)

func runRequest(B *testing.B, r *gin.Engine, path string) {
	url := fmt.Sprintf(addrFormat, path)
	req := httptest.NewRequest("GET", url, nil)

	// 创建一个假的writer
	w := httptest.NewRecorder()

	B.ReportAllocs()
	B.ResetTimer()

	for i := 0; i < B.N; i++ {
		r.ServeHTTP(w, req)
	}
}

func runHandle(B *testing.B, handler gin.HandlerFunc) {
	req := httptest.NewRequest("GET", "http://localhost/foo", nil)
	c := &gin.Context{
		Req:    req,
		Writer: gin.NewResponseWriter(httptest.NewRecorder(), 0, false),
		Engine: gin.New(),
	}

	B.ReportAllocs()
	B.ResetTimer()
	for i := 0; i < B.N; i++ {
		c.SetIndex(int8(0))
		handler(c)
	}
}

func BenchmarkDefaultOnlyPing(B *testing.B) {
	engine := gin.New()

	engine.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	runRequest(B, engine, "ping")
}

func BenchmarkMiddlewareLogger(B *testing.B) {
	runHandle(B, gin.Logger(ioutil.Discard))
}
