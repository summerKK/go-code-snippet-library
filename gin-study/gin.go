package gin

import (
	"context"
	"html/template"
	"log"
	"math"
	"net"
	"net/http"
	"path"
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
/********** HandlerFunc *********/
/************************************/

type HandlerFunc func(c *Context)

/************************************/
/********** handlers404 *********/
/************************************/

type handlers404 struct {
	engine *Engine
}

func (h *handlers404) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handlers := h.engine.allHandlers(h.engine.handlers404...)
	c := h.engine.createContext(w, r, nil, handlers)
	c.Next()
	if !c.Writer.Written() {
		http.NotFound(w, r)
	}
	// 放回池子
	c.Engine.freeCtx(c)
}

/************************************/
/********** RouterGroup *********/
/************************************/

// 路由
type RouterGroup struct {
	// 中间件
	Handlers []HandlerFunc
	// 前缀
	prefix string
	parent *RouterGroup
	engine *Engine
}

// 创建context,贯穿整个请求
func (r *RouterGroup) createContext(w http.ResponseWriter, req *http.Request, params httprouter.Params, handlers []HandlerFunc) *Context {
	ctx := r.engine.ctxPool.Get().(*Context)
	// 初始化ctx
	ctx.Writer.Reset(w)
	ctx.Request = req
	ctx.Params = params
	ctx.handlers = handlers
	ctx.index = -1

	return ctx
}

// 获取所有中间件
func (r *RouterGroup) allHandlers(handlers ...HandlerFunc) []HandlerFunc {
	local := append(r.Handlers, handlers...)
	if r.parent != nil {
		// 获取parent的中间件
		return r.parent.allHandlers(local...)
	}

	return local
}

//  添加中间件
func (r *RouterGroup) Use(middlewares ...HandlerFunc) {
	r.Handlers = append(r.Handlers, middlewares...)
}

// 返回新的group
func (r *RouterGroup) Group(component string, handlers ...HandlerFunc) *RouterGroup {
	prefix := r.pathFor(component)
	return &RouterGroup{
		Handlers: handlers,
		prefix:   prefix,
		parent:   r,
		engine:   r.engine,
	}
}

func (r *RouterGroup) pathFor(p string) string {
	joined := path.Join(r.prefix, p)
	if len(p) > 0 && p[len(p)-1] == '/' && joined[len(joined)-1] != '/' {
		joined += "/"
	}

	return joined
}

func (r *RouterGroup) Handle(method, p string, handlers []HandlerFunc) {
	pathName := path.Join(r.prefix, p)
	// 获取所有的中间件
	allHandlers := r.allHandlers(handlers...)
	// 处理请求
	r.engine.router.Handle(method, pathName, func(writer http.ResponseWriter, req *http.Request, params httprouter.Params) {
		// 创建context
		ctx := r.createContext(writer, req, params, allHandlers)
		ctx.Next()
		// 回收ctx
		r.engine.freeCtx(ctx)
	})
}

func (r *RouterGroup) GET(path string, handlers ...HandlerFunc) {
	r.Handle("GET", path, handlers)
}

func (r *RouterGroup) POST(path string, handlers ...HandlerFunc) {
	r.Handle("POST", path, handlers)
}

func (r *RouterGroup) DELETE(path string, handlers ...HandlerFunc) {
	r.Handle("DELETE", path, handlers)
}

func (r *RouterGroup) PATCH(path string, handlers ...HandlerFunc) {
	r.Handle("PATCH", path, handlers)
}

func (r *RouterGroup) PUT(path string, handlers ...HandlerFunc) {
	r.Handle("PUT", path, handlers)
}

func (r *RouterGroup) OPTIONS(path string, handlers ...HandlerFunc) {
	r.Handle("OPTIONS", path, handlers)
}

func (r *RouterGroup) Static(p, root string) {
	prefix := r.pathFor(p)
	p = path.Join(p, "/*filepath")
	// see https://studygolang.com/articles/9197
	fileServer := http.StripPrefix(prefix, http.FileServer(http.Dir(root)))
	r.GET(p, func(c *Context) {
		fileServer.ServeHTTP(c.Writer, c.Request)
	})
}

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
	handlers404 []HandlerFunc
	router      *httprouter.Router
	HTMLRender  render.Render
	// context pool
	ctxPool sync.Pool
	addr    string
}

func (e *Engine) LoadHTMLGlob(pattern string) {
	tmpl := template.Must(template.ParseGlob(pattern))
	e.SetHTMLTemplate(tmpl)
}

func (e *Engine) LoadHTMLFiles(files ...string) {
	tmpl := template.Must(template.ParseFiles(files...))
	e.SetHTMLTemplate(tmpl)
}

func (e *Engine) SetHTMLTemplate(tmpl *template.Template) {
	e.HTMLRender = render.HTMLRender{
		Template: tmpl,
	}
}

func (e *Engine) NotFound404(handler ...HandlerFunc) {
	e.handlers404 = handler
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
	e.run(c, addr, func(s *http.Server) error {
		return s.ListenAndServe()
	})
}

// https服务
func (e *Engine) RunTLS(c context.Context, addr string, cert string, key string) {
	e.run(c, addr, func(s *http.Server) error {
		return s.ListenAndServeTLS(cert, key)
	})
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
