package gin_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/matryer/is"
	"github.com/summerKK/go-code-snippet-library/gin-study"
)

var (
	url      = "http://localhost/%s"
	respText = "hello,world"
)

func init() {
	gin.SetMode(gin.TestModel)
}

func TestHandleStaticFile(t *testing.T) {
	assertIs := is.New(t)
	engine := gin.New()

	testRoot, _ := os.Getwd()
	f, err := ioutil.TempFile(testRoot, "")
	if err != nil {
		t.Fatal(err)
	}

	_, _ = f.WriteString("hello,world")
	f.Close()

	defer func() {
		_ = os.Remove(f.Name())
	}()

	fp := path.Join("/static", path.Base(f.Name()))
	req, err := http.NewRequest("GET", fp, nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	engine.Static("/static", "./")

	engine.ServeHTTP(w, req)

	assertIs.Equal(http.StatusOK, w.Code)
	assertIs.Equal("hello,world", w.Body.String())
	assertIs.Equal("text/plain; charset=utf-8", w.Header().Get("Content-Type"))
}

func TestHandleStaticDir(t *testing.T) {
	assertIs := is.New(t)
	engine := gin.New()

	engine.Static("/static", "./")

	req := httptest.NewRequest("GET", "/static/gin.go", nil)
	w := httptest.NewRecorder()

	engine.ServeHTTP(w, req)

	//assertIs.Equal(http.StatusOK, w.Code)

	bodyAsString := w.Body.String()
	assertIs.True(strings.Contains(bodyAsString, "package gin"))
	assertIs.Equal("text/plain; charset=utf-8", w.Header().Get("Content-Type"))
}

func TestEngine_Run(t *testing.T) {
	assertIs := is.New(t)
	engine := gin.New()

	engine.GET("/", func(c *gin.Context) {
		c.Abort(http.StatusOK)
	})

	engine.Run(context.Background(), ":8899")

	resp, err := http.Get("http://127.0.0.1:8899")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	assertIs.Equal(http.StatusOK, resp.StatusCode)
}

func TestEngine_LoadHTMLGlob(t *testing.T) {
	engine := gin.New()
	// 输出所有错误
	engine.Use(gin.ErrorLogger())

	assertIs := is.New(t)

	type templateFile struct {
		name    string
		content string
	}

	createTestDir := func(files []templateFile) string {
		dir, err := ioutil.TempDir("", "template")
		if err != nil {
			t.Fatal(err)
		}
		for _, file := range files {
			f, err := os.Create(filepath.Join(dir, file.name))
			if err != nil {
				t.Fatal(err)
			}
			_, err = f.Write([]byte(file.content))
			if err != nil {
				t.Fatal(err)
			}
			f.Close()
		}

		return dir
	}

	dir := createTestDir([]templateFile{
		// T0.tmpl is a plain template file that just invokes T1.
		{"T0.tmpl", `{{.summer}}`},
		// T1.tmpl defines a template, T1 that invokes T2.
		{"T1.tmpl", `{{define "T1"}}T1 invokes T2: ({{template "T2"}}){{end}}`},
		// T2.tmpl defines a template T2.
		{"T2.tmpl", `{{define "T2"}}This is T2{{end}}`},
	})

	defer os.RemoveAll(dir)

	pattern := filepath.Join(dir, "*.tmpl")

	engine.LoadHTMLGlob(pattern)

	engine.GET("/template", func(c *gin.Context) {
		c.HTML(http.StatusOK, "T0.tmpl", map[string]string{"summer": respText})
	})

	req := httptest.NewRequest("GET", fmt.Sprintf(url, "template"), nil)
	w := httptest.NewRecorder()

	engine.ServeHTTP(w, req)

	assertIs.Equal(http.StatusOK, w.Code)
	assertIs.Equal(respText, w.Body.String())
}
