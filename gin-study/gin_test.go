package gin_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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

func TestContext_ParseBody(t *testing.T) {
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
		err := c.ParseBody(&params0)
		if err != nil {
			t.Errorf("ParseBody got error:%v", err)
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
	basicAccounts := []gin.Account{
		{
			User:     u,
			Password: p,
		},
	}
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
