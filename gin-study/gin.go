package gin

import (
	"context"
	"html/template"
	"log"
	"math"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/summerKK/go-code-snippet-library/gin-study/render"
)

const (
	AbortIndex         = math.MaxInt8 / 2
	DefaultCtxPoolSize = 1024
)

/************************************/
/********** Engine *********/
/************************************/

type Config struct {
	CtxPoolSize    int
	CtxPreloadSize int
}

// 整个framework
type Engine struct {
	*RouterGroup
	// api未找到,触发的方法
	noRoute []HandlerFunc
	// 路由未找到触发的handle,`finalNoRoute`包含了`noRoute`
	finalNoRoute []HandlerFunc
	router       *httprouter.Router
	HTMLRender   render.Render
	// context pool
	ctxPool sync.Pool
	addr    string
}

func (e *Engine) LoadHTMLGlob(pattern string) {
	if IsDebugging() {
		render.HTMLDebug.AddGlob(pattern)
		e.HTMLRender = render.HTMLDebug
	} else {
		tmpl := template.Must(template.ParseGlob(pattern))
		e.SetHTMLTemplate(tmpl)
	}
}

func (e *Engine) LoadHTMLFiles(files ...string) {
	if IsDebugging() {
		render.HTMLDebug.AddFiles(files...)
		e.HTMLRender = render.HTMLDebug
	} else {
		tmpl := template.Must(template.ParseFiles(files...))
		e.SetHTMLTemplate(tmpl)
	}
}

func (e *Engine) SetHTMLTemplate(tmpl *template.Template) {
	e.HTMLRender = render.HTMLRender{
		Template: tmpl,
	}
}

func (e *Engine) NoRoute(handler ...HandlerFunc) {
	e.noRoute = handler
	e.finalNoRoute = e.combineHandlers(e.noRoute)
}

// ServeHTTP makes the router implement the http.Handler interface.
func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	e.router.ServeHTTP(w, req)
}

// 检查服务是否已开启
func (e *Engine) checkServerIsRunning(wg *sync.WaitGroup) {
	go func() {
		tick := time.Tick(time.Millisecond * 500)
		for {
			select {
			case <-tick:
				conn, err := net.DialTimeout("tcp", e.addr, time.Millisecond*500)
				// 服务已启动,记得return,停掉goroutine
				if err == nil {
					_ = conn.Close()
					wg.Done()
					return
				}
			}
		}
	}()
}

func (e *Engine) run(c context.Context, addr string, handle func(s *http.Server) error) {
	var wg sync.WaitGroup
	wg.Add(1)
	e.addr = addr

	go func() {
		s := &http.Server{
			Addr:    addr,
			Handler: e,
		}

		e.checkServerIsRunning(&wg)

		go func() {
			err := handle(s)
			if err != nil && err != http.ErrServerClosed {
				log.Fatalf("http server start error:%v\n", err)
			}
		}()

		<-c.Done()

		_ = s.Shutdown(context.Background())

		log.Printf("http server stopped!\n")
	}()

	wg.Wait()
}

// http服务
func (e *Engine) Run(c context.Context, addr string) {
	if IsDebugging() {
		debugPrint("Listening and serving HTTP on " + addr)
	}

	e.run(c, addr, func(s *http.Server) error {
		return s.ListenAndServe()
	})
}

// https服务
func (e *Engine) RunTLS(c context.Context, addr string, cert string, key string) {
	if IsDebugging() {
		debugPrint("Listening and serving HTTPS on " + addr)
	}

	e.run(c, addr, func(s *http.Server) error {
		return s.ListenAndServeTLS(cert, key)
	})
}

// 创建context,贯穿整个请求
func (e *Engine) createContext(w http.ResponseWriter, req *http.Request, params httprouter.Params, handlers []HandlerFunc) *Context {
	ctx := e.ctxPool.Get().(*Context)
	// 初始化ctx
	ctx.Writer.Reset(w)
	ctx.Request = req
	ctx.Params = params
	ctx.handlers = handlers
	ctx.index = -1
	// 情况错误
	ctx.Errors = ctx.Errors[0:0]
	ctx.accepted = nil

	return ctx
}

// 放回池子
func (e *Engine) freeCtx(c *Context) {
	e.ctxPool.Put(c)
}

func New() *Engine {
	engine := &Engine{}
	engine.RouterGroup = &RouterGroup{
		nil,
		"",
		nil,
		engine,
	}
	engine.router = httprouter.New()
	engine.router.NotFound = &handlers404{engine: engine}
	engine.ctxPool.New = func() interface{} {
		return &Context{
			Engine: engine,
			Writer: &ResponseWriter{},
		}
	}

	return engine
}

func Default() *Engine {
	engine := New()
	engine.Use(Recovery(), Logger(nil))

	return engine
}
