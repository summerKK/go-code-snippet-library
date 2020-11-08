package gin_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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

func TestMain(m *testing.M) {
	engine = gin.Default()
	ctx, cancelFunc = context.WithCancel(context.Background())
	go func() {
		engine.Run(ctx, fmt.Sprintf(":%s", port))
	}()
	time.Sleep(time.Second)
	m.Run()
}

func TestEngine_Run(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		time.Sleep(time.Second)
		cancelFunc()
		wg.Done()
	}()

	wg.Wait()
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
	resp, err := http.Post(fmt.Sprintf("http://127.0.0.1:%s/userinfo", port), "application/json", bytes.NewReader(jsonBytes))
	if err != nil {
		t.Fatal(err)
	}

	resp.Body.Close()
}
