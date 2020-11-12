package gin

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"path"

	"github.com/julienschmidt/httprouter"
	"github.com/summerKK/go-code-snippet-library/gin-study/binding"
)

const (
	AbortIndex         = math.MaxInt8 / 2
	DefaultCtxPoolSize = 1024
	MIMEJSON           = "application/json"
	MIMEHTML           = "text/html"
	MIMEXML            = "application/xml"
	MIMEXML2           = "text/xml"
	MIMEPlain          = "text/plain"
)

/************************************/
/********** ErrorMsg *********/
/************************************/

type ErrorMsg struct {
	Err  string      `json:"error"`
	Meta interface{} `json:"meta"`
}

type ErrorMsgs []ErrorMsg

func (e ErrorMsgs) String() string {
	var buf bytes.Buffer
	for i, msg := range e {
		text := fmt.Sprintf("Error #%02d: %s\n     Meta:%v\n\n", i+1, msg.Err, msg.Meta)
		buf.WriteString(text)
	}
	buf.WriteByte('\n')
	return buf.String()
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
	// 放回池子
	c.engine.reuseCtx(c)
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
	select {
	case ctx := <-r.engine.ctxPool:
		// 重置
		ctx.Writer.Reset(w)
		ctx.Req = req
		ctx.Params = params
		ctx.handlers = handlers
		ctx.index = -1
		return ctx
	default:
		return &Context{
			Req: req,
			Writer: &ResponseWriter{
				ResponseWriter: w,
			},
			index:    -1,
			engine:   r.engine,
			Params:   params,
			handlers: handlers,
			keep:     false,
		}
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
		ctx := r.createContext(writer, req, params, allHandlers)
		ctx.Next()
		// 回收ctx
		r.engine.reuseCtx(ctx)
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

/************************************/
/********** Engine *********/
/************************************/

type Config struct {
	CtxPoolSize    int
	CtxPreloadSize int
}

func (c Config) Check() error {
	if c.CtxPoolSize < 2 {
		return errors.New("CtxPoolSize must be at least 2")
	}

	if c.CtxPreloadSize > c.CtxPoolSize {
		return errors.New("CtxPreloadSize must be less or equal to CtxPoolSize")
	}

	return nil
}

// 整个framework
type Engine struct {
	*RouterGroup
	// api未找到,触发的方法
	handlers404   []HandlerFunc
	router        *httprouter.Router
	HTMLTemplates *template.Template
	// context pool
	ctxPool chan *Context
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

// 放回池子
func (e *Engine) reuseCtx(c *Context) {
	if c.keep {
		select {
		case e.ctxPool <- c:
		default:
		}
	}
}

func New() *Engine {
	return NewWithConfig(Config{
		CtxPoolSize:    DefaultCtxPoolSize,
		CtxPreloadSize: DefaultCtxPoolSize / 2,
	})
}

func NewWithConfig(c Config) *Engine {
	// 检查参数
	if err := c.Check(); err != nil {
		panic(err)
	}

	engine := &Engine{}
	engine.RouterGroup = &RouterGroup{
		nil,
		"/",
		nil,
		engine,
	}
	engine.router = httprouter.New()
	engine.router.NotFound = &handlers404{engine: engine}
	engine.ctxPool = make(chan *Context, c.CtxPoolSize)

	// 初始化ctx池
	for i := 0; i < c.CtxPreloadSize; i++ {
		engine.ctxPool <- &Context{
			keep:   true,
			engine: engine,
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

/************************************/
/********** Context *********/
/************************************/

// gin的核心模块.
type Context struct {
	ID     int
	Req    *http.Request
	Writer ResponseWriterInterface
	Keys   map[string]interface{}
	Params httprouter.Params
	// 收集错误.在logger中间件进行记录
	Errors ErrorMsgs
	// 中间件
	handlers []HandlerFunc
	index    int8
	engine   *Engine
	// 标识是否需要放回池子中
	keep bool
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

//  标识ctx可以放入池子中
func (c *Context) Keep() {
	// 修改keep状态.使其可以放入池子中继续使用
	if !c.keep {
		c.keep = true
	} else {
		log.Println("gin: trying to Keep same context several times")
	}
}

// 释放ctx
func (c *Context) Release() {
	if c.keep {
		c.keep = false
	} else {
		log.Println("gin: trying to release same context several times")
	}
}

// 终止请求
func (c *Context) Abort(code int) {
	if code >= 0 {
		c.Writer.WriteHeader(code)
	}
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
		Err:  err.Error(),
		Meta: meta,
	})
}

func (c *Context) LastError() error {
	l := len(c.Errors)
	if l > 0 {
		return errors.New(c.Errors[l-1].Err)
	}

	return nil
}

func (c *Context) Set(key string, v interface{}) {
	if c.Keys == nil {
		c.Keys = make(map[string]interface{})
	}

	c.Keys[key] = v
}

func (c *Context) Get(key string) (interface{}, error) {
	if c.Keys != nil {
		item, ok := c.Keys[key]
		if ok {
			return item, nil
		}
	}

	return nil, errors.New(fmt.Sprintf("Key %s does'nt exist", key))
}

func (c *Context) MustGet(key string) interface{} {
	v, err := c.Get(key)
	if v == nil || err != nil {
		log.Panicf("MustGet key %s does'nt exist", key)
	}

	return v
}

func (c *Context) filterFlags(content string) string {
	for i, c := range content {
		if c == ' ' || c == ';' {
			return content[:i]
		}
	}

	return content
}

func (c *Context) Bind(v interface{}) bool {
	var b binding.Binding
	contentType := c.filterFlags(c.Req.Header.Get("Content-Type"))
	switch contentType {
	case MIMEJSON:
		b = binding.JSON
	case MIMEXML, MIMEXML2:
		b = binding.XML
	default:
		c.Fail(400, errors.New("unknown content-type: "+contentType))
		return false
	}

	return c.BindWith(v, b)
}

func (c *Context) BindWith(v interface{}, b binding.Binding) bool {
	if err := b.Bind(c.Req, v); err == nil {
		return true
	}

	return false
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

func (c *Context) SetIndex(index int8) {
	c.index = index
}
