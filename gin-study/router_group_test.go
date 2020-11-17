package gin_test

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/matryer/is"
	"github.com/summerKK/go-code-snippet-library/gin-study"
)

func TestRouterGroup_Use(t *testing.T) {
	assertIs := is.New(t)
	engine := gin.New()
	engine.Use(func(c *gin.Context) {
		log.Println("<<<<<<<<<<<<<")
		c.Next()
		log.Println("<<<<<<<<<<<<<")
	})

	group := engine.Group("/api")

	group.Use(func(c *gin.Context) {
		log.Println("      >>>>>>>>>>>>>")
		c.Next()
		log.Println("      >>>>>>>>>>>>>")
	})

	group.Use(func(c *gin.Context) {
		log.Println("            >>>>>>>>>>>>>")
		c.Next()
		log.Println("            >>>>>>>>>>>>>")
	})

	group.GET("/middleware", func(c *gin.Context) {
		_, _ = c.Writer.Write([]byte(respText))
		log.Println("                   hello,world")
		c.Abort(200)
	})

	req := httptest.NewRequest("GET", fmt.Sprintf(url, "api/middleware"), nil)
	w := httptest.NewRecorder()

	engine.ServeHTTP(w, req)

	assertIs.Equal(http.StatusOK, w.Code)
	assertIs.Equal(respText, w.Body.String())
}

func TestRouterGroup_Group(t *testing.T) {
	assertIs := is.New(t)
	engine := gin.New()

	r := engine.Group("/api", func(c *gin.Context) {
		c.Next()
		_, _ = c.Writer.Write([]byte(respText))
	})

	r.GET("/test", func(c *gin.Context) {

	})

	req := httptest.NewRequest("GET", fmt.Sprintf(url, "api/test"), nil)
	w := httptest.NewRecorder()

	engine.ServeHTTP(w, req)

	assertIs.Equal(respText, w.Body.String())
}

func TestRouterGroupParams(t *testing.T) {
	assertIs := is.New(t)
	engine := gin.New()
	username := "summer"

	engine.GET("/user/:username", func(c *gin.Context) {
		assertIs.Equal(c.Params.ByName("username"), username)
	})

	req := httptest.NewRequest("GET", fmt.Sprintf(url, "user/"+username), nil)
	w := httptest.NewRecorder()

	engine.ServeHTTP(w, req)
}
