package gin_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/matryer/is"
	"github.com/summerKK/go-code-snippet-library/gin-study"
)

func TestBasicAuth(t *testing.T) {
	assertIs := is.New(t)
	engine := gin.New()

	u := "summer"
	p := "123"
	var basicAccounts gin.Accounts = map[string]string{}
	basicAccounts[u] = p
	engine.Use(gin.BasicAuth(basicAccounts))

	engine.GET("/basic-auth", func(c *gin.Context) {
		c.String(http.StatusOK, respText)
	})

	req, err := http.NewRequest("GET", fmt.Sprintf(url, "basic-auth"), nil)
	if err != nil {
		t.Fatal(err)
	}

	req.SetBasicAuth(u, p)

	w := httptest.NewRecorder()

	engine.ServeHTTP(w, req)

	assertIs.Equal(http.StatusOK, w.Code)
	assertIs.Equal(respText, w.Body.String())

	// 401 test
	req = httptest.NewRequest("GET", fmt.Sprintf(url, "basic-auth"), nil)
	req.SetBasicAuth("sunny", "dddd")
	w = httptest.NewRecorder()

	engine.ServeHTTP(w, req)

	assertIs.Equal(http.StatusUnauthorized, w.Code)
}
