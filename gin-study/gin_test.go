package gin_test

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/matryer/is"
	"github.com/summerKK/go-code-snippet-library/gin-study"
)

var engine *gin.Engine
var ctx context.Context
var cancelFunc func()
var port string = "8080"
var addrFormat = "http://127.0.0.1:" + port + "/%s"
var respText = "hello,world"

func TestMain(m *testing.M) {
	engine = gin.Default()
	engine.Use(gin.ErrorLogger())
	ctx, cancelFunc = context.WithCancel(context.Background())
	defer cancelFunc()

	go func() {
		engine.Run(ctx, fmt.Sprintf(":%s", port))
	}()
	time.Sleep(time.Second)
	m.Run()

	cancelFunc()
}

func TestRouterGroup_Use(t *testing.T) {
	assertIs := is.New(t)
	engine.Use(func(c *gin.Context) {
		c.Next()
		_, _ = c.Writer.Write([]byte(respText))
	})

	engine.GET("/middleware", func(c *gin.Context) {
		c.Abort(200)
	})

	resp, err := http.Get(fmt.Sprintf(addrFormat, "middleware"))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	all, _ := ioutil.ReadAll(resp.Body)
	assertIs.Equal(string(all), respText)
}

func TestContext_Bind(t *testing.T) {
	assertIs := is.New(t)

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
		r := c.Bind(&params0)
		if !r {
			t.Error("Bind params got error")
			return
		}

		assertIs.Equal(params0.Name, params1.Name)
		assertIs.Equal(params0.Age, params1.Age)
		c.Abort(200)
	})

	values := make(map[string]interface{})
	values["Name"] = params1.Name
	values["Age"] = params1.Age
	jsonBytes, _ := json.Marshal(values)
	resp, err := http.Post(fmt.Sprintf(addrFormat, "userinfo"), "application/json", bytes.NewReader(jsonBytes))
	if err != nil {
		t.Fatal(err)
	}

	resp.Body.Close()
}

func TestContext_Bind2(t *testing.T) {
	assertIs := is.New(t)

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

	resp, err := http.Get(fmt.Sprintf(addrFormat, "userinfo?name=summer&age=28&name=summer&age=28"))
	if err != nil {
		t.Fatal(err)
	}

	resp.Body.Close()
}

func TestRouterGroup_Group(t *testing.T) {
	assertIs := is.New(t)
	r := engine.Group("/api", func(c *gin.Context) {
		c.Next()
		_, _ = c.Writer.Write([]byte(respText))
	})

	r.GET("/test", func(c *gin.Context) {

	})

	resp, err := http.Get(fmt.Sprintf(addrFormat, "api/test"))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	all, _ := ioutil.ReadAll(resp.Body)
	assertIs.Equal(string(all), respText)
}

func TestRouterGroupParams(t *testing.T) {
	assertIs := is.New(t)
	username := "summer"
	engine.GET("/user/:username", func(c *gin.Context) {
		assertIs.Equal(c.Params.ByName("username"), username)
	})

	resp, err := http.Get(fmt.Sprintf(addrFormat, "/user/"+username))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
}

func TestBasicAuth(t *testing.T) {
	assertIs := is.New(t)
	u := "summer"
	p := "123"
	var basicAccounts gin.Accounts = map[string]string{}
	basicAccounts[u] = p
	engine.Use(gin.BasicAuth(basicAccounts))

	engine.GET("/basic-auth", func(c *gin.Context) {
		c.String(http.StatusOK, respText)
	})

	req, err := http.NewRequest("GET", fmt.Sprintf(addrFormat, "basic-auth"), nil)
	if err != nil {
		t.Fatal(err)
	}

	req.SetBasicAuth(u, p)
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()

	all, _ := ioutil.ReadAll(response.Body)

	assertIs.Equal(string(all), respText)
}

func TestErrorMsgs_String(t *testing.T) {
	var errorMsgs gin.ErrorMsgs = []gin.ErrorMsg{
		{
			Err:  "hello,world",
			Meta: map[string]int{"a": 1},
		},
		{
			Err:  "hello,world",
			Meta: []int{1, 2},
		},
	}

	fmt.Println(errorMsgs.String())
}

func TestContext_Pool(t *testing.T) {
	engine.GET("/userinfo", func(c *gin.Context) {
		c.Abort(http.StatusOK)
	})

	wg := &sync.WaitGroup{}
	for i := 0; i < gin.DefaultCtxPoolSize/8; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			resp, err := http.Get(fmt.Sprintf(addrFormat, "userinfo"))
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			wg.Done()
		}(wg)
	}

	wg.Wait()
}

func TestContext_JSON(t *testing.T) {
	assertIs := is.New(t)
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

	resp, err := http.Get(fmt.Sprintf(addrFormat, "json"))
	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()

	var p P
	decoder := json.NewDecoder(resp.Body)
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
