package gin_test

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/matryer/is"
	"github.com/summerKK/go-code-snippet-library/gin-study"
)

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

func TestContext_Redirect(t *testing.T) {
	engine := gin.New()
	assertIs := is.New(t)

	engine.GET("/redirect1", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/redirect0")
	})

	req := httptest.NewRequest("GET", fmt.Sprintf(url, "redirect1"), nil)
	w := httptest.NewRecorder()

	engine.ServeHTTP(w, req)

	assertIs.Equal(http.StatusFound, w.Code)
}
