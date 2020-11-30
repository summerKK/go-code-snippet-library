package gin

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/summerKK/go-code-snippet-library/gin-study/binding"
	"github.com/summerKK/go-code-snippet-library/gin-study/render"
)

const (
	ErrorTypeInternal = 1 << iota
	ErrorTypeExternal = 1 << iota
	ErrorTypeAll      = (1 << 32) - 1
)

type errorMsg struct {
	Err  string      `json:"error"`
	Type uint32      `json:"-"`
	Meta interface{} `json:"meta"`
}

// 网络协商
type Negotiate struct {
	Offered  []string
	Data     interface{}
	JSONData interface{}
	XMLData  interface{}
	HTMLData interface{}
	HTMLPath string
}

type errorMsgs []errorMsg

// 返回特定的错误
func (e errorMsgs) ByType(typ uint32) errorMsgs {
	if len(e) == 0 {
		return e
	}

	result := make(errorMsgs, 0, len(e))
	for _, msg := range e {
		// 只返回特定的错误
		if msg.Type&typ > 0 {
			result = append(result, msg)
		}
	}

	return result
}

func (e errorMsgs) String() string {
	if len(e) == 0 {
		return ""
	}

	var buf bytes.Buffer
	for i, msg := range e {
		text := fmt.Sprintf("Error #%02d: %s\n     Meta:%v\n\n", i+1, msg.Err, msg.Meta)
		buf.WriteString(text)
	}
	buf.WriteByte('\n')
	return buf.String()
}

type H map[string]interface{}

func (h H) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name = xml.Name{Space: "", Local: "map"}

	if err := e.EncodeToken(start); err != nil {
		return err
	}

	for key, v := range h {
		elem := xml.StartElement{
			Name: xml.Name{Local: key},
			Attr: []xml.Attr{},
		}
		if err := e.EncodeElement(v, elem); err != nil {
			return err
		}
	}
	if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
		return err
	}

	return nil
}

// gin的核心模块.
type Context struct {
	Request *http.Request
	Writer  ResponseWriterInterface
	Keys    map[string]interface{}
	Params  httprouter.Params
	// 收集错误.在logger中间件进行记录
	Errors errorMsgs
	// 中间件
	handlers []HandlerFunc
	index    int8
	Engine   *Engine
	accepted []string
}

// 执行middleware
// gin的中间件实现很巧妙.主要是在 c.index < s.当c.index == 0 //标记为`a` 的时候.
// 在 c.handlers[0](c) 调用后.在c.handlers[0](c) {c.Next() //标记为`b`} 会调用下一次的Next()方法
// 此时在`a`的循环里面c.index已经变成了1
func (c *Context) Next() {
	// index 初始值为 -1
	c.index++
	// 避免数组越界.比如Abort的时候把index设置为AbortIndex
	if c.index < int8(len(c.handlers)) {
		c.handlers[c.index](c)
	}
}

func (c *Context) Abort(code int) {
	// 把index设置到最大,让剩余的中间件不执行
	c.index = AbortIndex
}

// 终止请求
func (c *Context) AbortWithStatus(code int) {
	c.Writer.WriteHeader(code)
	// 把index设置到最大,让剩余的中间件不执行
	c.index = AbortIndex
}

// 失败调用方法
func (c *Context) Fail(code int, err error) {
	c.Error(err, "Operation aborted")
	c.AbortWithStatus(code)
}

// 添加错误
func (c *Context) Error(err error, meta interface{}) {
	c.ErrorTyped(err, ErrorTypeExternal, meta)
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

func (c *Context) Bind(v interface{}) bool {
	var b binding.Binding
	contentType := filterFlags(c.Request.Header.Get("Content-Type"))
	switch {
	case c.Request.Method == "GET" || contentType == render.MIMEPOSTForm:
		b = binding.FORM
	case contentType == render.MIMEJSON:
		b = binding.JSON
	case contentType == render.MIMEXML || contentType == render.MIMEXML2:
		b = binding.XML
	default:
		c.Fail(400, errors.New("unknown content-type: "+contentType))
		return false
	}

	return c.BindWith(v, b)
}

func (c *Context) BindWith(v interface{}, b binding.Binding) bool {
	if err := b.Bind(c.Request, v); err == nil {
		return true
	}

	return false
}

func (c *Context) Render(render render.Render, code int, obj ...interface{}) {
	if err := render.Render(c.Writer, code, obj...); err != nil {
		c.ErrorTyped(err, ErrorTypeInternal, obj)
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func (c *Context) JSON(code int, v interface{}) {
	c.Render(render.JSON, code, v)
}

func (c *Context) XML(code int, v interface{}) {
	c.Render(render.XML, code, v)
}

// https://golang.org/pkg/text/template/#Template
func (c *Context) HTML(code int, name string, data interface{}) {
	c.Render(c.Engine.HTMLRender, code, name, data)
}

func (c *Context) String(code int, format string, args ...interface{}) {
	c.Render(render.Plain, code, format, args)
}

func (c *Context) Data(code int, contentType string, data []byte) {
	if len(contentType) > 0 {
		c.Writer.Header().Set("Content-Type", contentType)
	}

	c.Writer.WriteHeader(code)
	_, _ = c.Writer.Write(data)
}

func (c *Context) File(filepath string) {
	http.ServeFile(c.Writer, c.Request, filepath)
}

func (c *Context) Redirect(code int, location string) {
	if code >= 300 && code <= 308 {
		c.Render(render.Redirect, code, location)
	} else {
		panic(fmt.Sprintf("Cannot send a redirect with status code %d", code))
	}
}

func (c *Context) SetIndex(index int8) {
	c.index = index
}

func (c *Context) ErrorTyped(err error, typ uint32, meta interface{}) {
	c.Errors = append(c.Errors, errorMsg{
		Err:  err.Error(),
		Type: typ,
		Meta: meta,
	})
}

func (c *Context) Negotiate(code int, config Negotiate) {
	switch c.NegotiateFormat(config.Offered...) {
	case render.MIMEJSON:
		data := chooseData(config.JSONData, config.Data)
		c.JSON(code, data)

	case render.MIMEHTML:
		data := chooseData(config.HTMLData, config.Data)
		if len(config.HTMLPath) == 0 {
			panic("negotiate config is wrong. html path is needed.")
		}
		c.HTML(code, config.HTMLPath, data)

	case render.MIMEXML:
		data := chooseData(config.XMLData, config.Data)
		c.XML(code, data)

	default:
		c.Fail(http.StatusNotAcceptable, errors.New("the accepted formats are not offered by the server"))
	}
}

func (c *Context) NegotiateFormat(offered ...string) string {
	if len(offered) == 0 {
		panic("you must provide at least one offer")
	}

	if c.accepted == nil {
		c.accepted = parseAccept(c.Request.Header.Get("Accept"))
	}

	if len(c.accepted) == 0 {
		return offered[0]
	} else {
		for _, accepted := range c.accepted {
			for _, offert := range offered {
				if accepted == offert {
					return offert
				}
			}
		}
	}

	return ""
}

func (c *Context) SetAccepted(formats ...string) {
	c.accepted = formats
}
