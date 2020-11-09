package gin

import (
	"context"
	"encoding/json"
	"html/template"
	"log"
	"math"
	"net/http"
	"path"

	"github.com/julienschmidt/httprouter"
)

const (
	AbortIndex = math.MaxInt8 / 2
)

/************************************/
/********** ErrorMsg *********/
/************************************/

type ErrorMsg struct {
	Message string      `json:"msg"`
	Meta    interface{} `json:"meta"`
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
/********** responseWriter *********/
/************************************/
type ResponseWriter interface {
	http.ResponseWriter
	Status() int
	Written() bool
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (w *responseWriter) WriteHeader(s int) {
	w.ResponseWriter.WriteHeader(s)
	w.status = s
}

func (w *responseWriter) Write(b []byte) (int, error) {
	return w.ResponseWriter.Write(b)
}

func (w *responseWriter) Status() int {
	return w.status
}

// 判断是否已经写入数据
func (w *responseWriter) Written() bool {
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
		Writer:   &responseWriter{w, 0},
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

// 返回新的group
func (r *RouterGroup) Group(component string, handlers ...HandlerFunc) *RouterGroup {
	prefix := path.Join(r.prefix, component)
	return &RouterGroup{
		Handlers: handlers,
		prefix:   prefix,
		parent:   r,
		engine:   r.engine,
	}
}

func (r *RouterGroup) Handle(method, p string, handlers []HandlerFunc) {
	pathName := path.Join(r.prefix, p)
	// 获取所有的中间件
	allHandlers := r.allHandlers(handlers...)
	// 处理请求
	r.engine.router.Handle(method, pathName, func(writer http.ResponseWriter, req *http.Request, params httprouter.Params) {
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
	handlers404   []HandlerFunc
	router        *httprouter.Router
	HTMLTemplates *template.Template
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

func (e *Engine) LoadHTMLTemplates(pattern string) {
	e.HTMLTemplates = template.Must(template.ParseGlob(pattern))
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
	Writer *responseWriter
	Keys   map[string]interface{}
	Params httprouter.Params
	// 收集错误.在logger中间件进行记录
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
	// 避免数组越界.比如Abort的时候把index设置为AbortIndex
	if c.index < int8(len(c.handlers)) {
		c.handlers[c.index](c)
	}
}

// 终止请求
func (c *Context) Abort(code int) {
	c.Writer.WriteHeader(code)
	// 把index设置到最大,让剩余的中间件不执行
	c.index = AbortIndex
}

// 失败调用方法
func (c *Context) Fail(code int, err error) {
	c.Error(err, "Operation aborted")
	c.Abort(code)
}

func (c *Context) Error(err error, meta interface{}) {
	c.Errors = append(c.Errors, ErrorMsg{
		Message: err.Error(),
		Meta:    meta,
	})
}

func (c *Context) Set(key string, v interface{}) {
	if c.Keys == nil {
		c.Keys = make(map[string]interface{})
	}

	c.Keys[key] = v
}

func (c *Context) Get(key string) interface{} {
	var ok bool
	if c.Keys == nil {
		ok = false
	}

	v, ok := c.Keys[key]

	if !ok {
		log.Panicf("Key %s does'nt exist", key)
	}

	return v
}

func (c *Context) EnsureBody(item interface{}) bool {
	if err := c.ParseBody(item); err != nil {
		c.Fail(400, err)
		return false
	}

	return true
}

// 解析请求参数
func (c *Context) ParseBody(item interface{}) error {
	decoder := json.NewDecoder(c.Req.Body)
	if err := decoder.Decode(&item); err == nil {
		return Validate(c, item)
	} else {
		return err
	}
}

func (c *Context) JSON(code int, v interface{}) {
	c.Writer.Header().Set("Content-Type", "application/json")
	if code >= 0 {
		c.Writer.WriteHeader(code)
	}
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(v); err != nil {
		c.Error(err, v)
		http.Error(c.Writer, err.Error(), 500)
	}
}

// https://golang.org/pkg/text/template/#Template
func (c *Context) HTML(code int, name string, data interface{}) {
	c.Writer.Header().Set("Content-Type", "text/html")
	if code >= 0 {
		c.Writer.WriteHeader(code)
	}
	if err := c.engine.HTMLTemplates.ExecuteTemplate(c.Writer, name, data); err != nil {
		c.Error(err, map[string]interface{}{
			"name": name,
			"data": data,
		})
		http.Error(c.Writer, err.Error(), 500)
	}
}

func (c *Context) String(code int, msg string) {
	if code >= 0 {
		c.Writer.WriteHeader(code)
	}

	c.Writer.Header().Set("Content-Type", "text/plain")
	_, _ = c.Writer.Write([]byte(msg))
}

func (c *Context) Data(code int, data []byte) {
	if code >= 0 {
		c.Writer.WriteHeader(code)
	}

	_, _ = c.Writer.Write(data)
}
