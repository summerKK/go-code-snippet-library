package gin

import (
	"fmt"
	"net/http"
	"path"

	"github.com/julienschmidt/httprouter"
	"github.com/summerKK/go-code-snippet-library/gin-study/render"
)

type HandlerFunc func(c *Context)

type handlers404 struct {
	engine *Engine
}

func (h *handlers404) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := h.engine.createContext(w, r, nil, h.engine.finalNoRoute)
	c.Writer.WriteHeader(http.StatusNotFound)
	c.Next()

	if !c.Writer.Written() {
		if c.Writer.Status() == http.StatusNotFound {
			c.Data(http.StatusNotFound, render.MIMEPlain, []byte("404 page not found"))
		} else {
			c.Writer.WriteHeaderNow()
		}
	}

	// 放回池子
	c.Engine.freeCtx(c)
}

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
	// 情况错误
	ctx.Errors = ctx.Errors[0:0]
	ctx.accepted = nil

	return ctx
}

//  添加中间件
func (r *RouterGroup) Use(middlewares ...HandlerFunc) {
	r.Handlers = append(r.Handlers, middlewares...)
	// 给找不到路由的handler赋值,并且放在最后执行
	r.engine.finalNoRoute = r.engine.combineHandlers(r.engine.noRoute)
}

// 返回新的group
func (r *RouterGroup) Group(component string, handlers ...HandlerFunc) *RouterGroup {
	prefix := r.pathFor(component)

	return &RouterGroup{
		Handlers: r.combineHandlers(handlers),
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
	pathName := r.pathFor(p)
	// 获取所有的中间件
	combinedHandlers := r.combineHandlers(handlers)
	if ginMode == debugCode {
		numHandlers := len(combinedHandlers)
		name := funcName(combinedHandlers[numHandlers-1])
		fmt.Printf("[GIN-debug] %-5s %-25s --> %s (%d handlers)\n", method, p, name, numHandlers)
	}
	// 处理请求
	r.engine.router.Handle(method, pathName, func(writer http.ResponseWriter, req *http.Request, params httprouter.Params) {
		// 创建context
		ctx := r.createContext(writer, req, params, combinedHandlers)
		ctx.Next()
		// 添加响应头
		ctx.Writer.WriteHeaderNow()
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

func (r *RouterGroup) HEAD(path string, handlers ...HandlerFunc) {
	r.Handle("HEAD", path, handlers)
}

func (r *RouterGroup) Static(p, root string) {
	prefix := r.pathFor(p)
	p = path.Join(p, "/*filepath")
	// see https://studygolang.com/articles/9197
	fileServer := http.StripPrefix(prefix, http.FileServer(http.Dir(root)))
	r.GET(p, func(c *Context) {
		fileServer.ServeHTTP(c.Writer, c.Request)
	})

	r.HEAD(p, func(c *Context) {
		fileServer.ServeHTTP(c.Writer, c.Request)
	})
}

func (r *RouterGroup) combineHandlers(handlers []HandlerFunc) []HandlerFunc {
	if len(handlers) == 0 {
		return r.Handlers
	}

	l := len(r.Handlers) + len(handlers)
	h := make([]HandlerFunc, 0, l)
	h = append(h, r.Handlers...)
	h = append(h, handlers...)

	return h
}
