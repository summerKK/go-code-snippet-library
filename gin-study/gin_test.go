package gin_test

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/matryer/is"
	"github.com/summerKK/go-code-snippet-library/gin-study"
)

var (
	url      = "http://localhost/%s"
	respText = "hello,world"
)

func TestRouterGroup_Use(t *testing.T) {
	assertIs := is.New(t)
	engine := gin.New()

	engine.Use(func(c *gin.Context) {
		c.Next()
		_, _ = c.Writer.Write([]byte(respText))
	})

	engine.GET("/middleware", func(c *gin.Context) {
		c.Abort(200)
	})

	req := httptest.NewRequest("GET", fmt.Sprintf(url, "middleware"), nil)
	w := httptest.NewRecorder()

	engine.ServeHTTP(w, req)

	assertIs.Equal(http.StatusOK, w.Code)
	assertIs.Equal(respText, w.Body.String())
}

func TestContext_Bind(t *testing.T) {
	assertIs := is.New(t)
	engine := gin.New()

	type P struct {
		Name string `json:"name" binding:"required"`
		Age  int    `json:"age" binding:"required"`
	}

	params1 := P{
		Name: "summer",
		Age:  28,
	}

	engine.POST("/userinfo", func(c *gin.Context) {
		var params0 P

		assertIs.True(c.Bind(&params0))
		assertIs.Equal(params0.Name, params1.Name)
		assertIs.Equal(params0.Age, params1.Age)
		c.Abort(200)
	})

	values := make(map[string]interface{})
	values["Name"] = params1.Name
	values["Age"] = params1.Age
	jsonBytes, _ := json.Marshal(values)

	req := httptest.NewRequest("POST", fmt.Sprintf(url, "userinfo"), bytes.NewReader(jsonBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	engine.ServeHTTP(w, req)
}

func TestContext_Bind2(t *testing.T) {
	assertIs := is.New(t)
	engine := gin.New()

	type P struct {
		Name []string `form:"name" binding:"required"`
		Age  []int    `form:"age" binding:"required"`
	}

	engine.GET("/userinfo", func(c *gin.Context) {
		var params0 P
		if c.Bind(&params0) {
			assertIs.Equal([]string{"summer", "summer"}, params0.Name)
			assertIs.Equal([]int{28, 28}, params0.Age)
			c.Abort(http.StatusOK)
			return
		}

		t.Error("Bind params got error")
	})

	req := httptest.NewRequest("GET", fmt.Sprintf(url, "userinfo?name=summer&age=28&name=summer&age=28"), nil)
	w := httptest.NewRecorder()

	engine.ServeHTTP(w, req)
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

	assertIs.Equal(respText, w.Body.String())
}

func TestContext_Pool(t *testing.T) {
	assertIs := is.New(t)
	engine := gin.New()

	engine.GET("/pool", func(c *gin.Context) {
		c.Abort(http.StatusOK)
	})

	wg := &sync.WaitGroup{}
	for i := 0; i < gin.DefaultCtxPoolSize; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			req := httptest.NewRequest("GET", fmt.Sprintf(url, "pool"), nil)
			w := httptest.NewRecorder()

			engine.ServeHTTP(w, req)

			assertIs.Equal(http.StatusOK, w.Code)

			wg.Done()
		}(wg)
	}

	wg.Wait()
}

func TestContext_JSON(t *testing.T) {
	assertIs := is.New(t)
	engine := gin.New()

	type P struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	engine.GET("/json", func(c *gin.Context) {
		c.JSON(http.StatusOK, P{
			Name: "summer",
			Age:  28,
		})
	})

	req := httptest.NewRequest("GET", fmt.Sprintf(url, "json"), nil)
	w := httptest.NewRecorder()

	engine.ServeHTTP(w, req)

	var p P
	decoder := json.NewDecoder(w.Body)
	if err := decoder.Decode(&p); err != nil {
		t.Fatal(err)
	}

	assertIs.Equal(p.Name, "summer")
	assertIs.Equal(p.Age, 28)
}

func TestH_MarshalXML(t *testing.T) {
	h := gin.H{
		"slice": []int{1, 2, 3},
		"name":  "summer",
	}

	b, err := xml.Marshal(h)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(b))
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
