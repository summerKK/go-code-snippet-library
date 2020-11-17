package gin_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/matryer/is"
	"github.com/summerKK/go-code-snippet-library/gin-study"
)

func TestPanicInHandler(t *testing.T) {
	engine := gin.Default()
	assertIs := is.New(t)

	engine.GET("/panic", func(c *gin.Context) {
		panic("hello,world")
	})

	req := httptest.NewRequest("GET", fmt.Sprintf(url, "panic"), nil)
	w := httptest.NewRecorder()

	engine.ServeHTTP(w, req)

	assertIs.Equal(http.StatusInternalServerError, w.Code)
}
