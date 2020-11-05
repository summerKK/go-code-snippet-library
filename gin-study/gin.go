package gin

import (
	"context"
	"log"
	"net/http"
	"path"
)
import "github.com/julienschmidt/httprouter"

/************************************/
/********** ErrorMsg *********/
/************************************/

type ErrorMsg struct {
	Message string `json:"msg"`
	Meta    string `json:"meta"`
}

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
}

/************************************/
/********** ResponseWriter *********/
/************************************/

type ResponseWriter struct {
	http.ResponseWriter
	status int
}

func (w *ResponseWriter) WriteHeader(s int) {
	w.ResponseWriter.WriteHeader(s)
	w.status = s
}

func (w *ResponseWriter) Write(b []byte) (int, error) {
	return w.ResponseWriter.Write(b)
}

func (w *ResponseWriter) Status() int {
	return w.status
}

// 判断是否已经写入数据
func (w *ResponseWriter) Written() bool {
	return w.Status() != 0
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
	return &Context{
		Req:      req,
		Writer:   &ResponseWriter{w, 0},
		index:    -1,
		engine:   r.engine,
		Params:   params,
		handlers: handlers,
	}
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

func (r *RouterGroup) Handle(method, p string, handlers []HandlerFunc) {
	path := path.Join(r.prefix, p)
	allHandlers := r.allHandlers(handlers...)
	r.engine.router.Handle(method, path, func(writer http.ResponseWriter, req *http.Request, params httprouter.Params) {
		// 创建context
		r.createContext(writer, req, params, allHandlers).Next()
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

/************************************/
/********** Engine *********/
/************************************/

// 整个framework
type Engine struct {
	*RouterGroup
	// api未找到,触发的方法
	handlers404 []HandlerFunc
	router      *httprouter.Router
}

func (e *Engine) NotFound404(handler ...HandlerFunc) {
	e.handlers404 = handler
}

// ServeHTTP makes the router implement the http.Handler interface.
func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	e.router.ServeHTTP(w, req)
}

func (e *Engine) Run(c context.Context, addr string) {
	s := http.Server{
		Addr:    addr,
		Handler: e,
	}
	go func() {
		err := s.ListenAndServe()
		if err != http.ErrServerClosed {
			log.Fatalf("http server start error:%v\n", err)
		}
	}()

	<-c.Done()

	_ = s.Shutdown(context.Background())

	log.Printf("http server stopped!\n")
}

func New() *Engine {
	engine := &Engine{}
	engine.RouterGroup = &RouterGroup{
		nil,
		"/",
		nil,
		engine,
	}
	engine.router = httprouter.New()
	engine.router.NotFound = &handlers404{engine: engine}

	return engine
}

func Default() *Engine {
	engine := New()
	engine.Use(Recovery(), Logger())

	return engine
}

/************************************/
/********** Context *********/
/************************************/

// gin的核心模块.
type Context struct {
	Req    *http.Request
	Writer *ResponseWriter
	Keys   map[string]interface{}
	Params httprouter.Params
	Errors []ErrorMsg
	// 中间件
	handlers []HandlerFunc
	index    int8
	engine   *Engine
}

// 执行middleware
// gin的中间件实现很巧妙.主要是在 c.index < s.当c.index == 0 //标记为`a` 的时候.
// 在 c.handers[0](c) 调用后.在c.handers[0](c) {c.Next() //标记为`b`} 会调用下一次的Next()方法
// 此时在`a`的循环里面c.index已经变成了1
func (c *Context) Next() {
	// index 初始值为 -1
	c.index++
	_ = int8(len(c.handlers))
	c.handlers[c.index](c)
}
