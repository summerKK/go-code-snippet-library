package render

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html/template"
	"net/http"
)

const (
	MIMEJSON     = "application/json"
	MIMEHTML     = "text/html"
	MIMEXML      = "application/xml"
	MIMEXML2     = "text/xml"
	MIMEPlain    = "text/plain"
	MIMEPOSTForm = "application/x-www-form-urlencoded"
)

type Render interface {
	Render(http.ResponseWriter, int, ...interface{}) error
}

type jsonRender struct{}

type xmlRender struct{}

type plainRender struct{}

type redirectRender struct{}

type HTMLRender struct {
	Template *template.Template
}

type htmlDebugRender struct {
	globs []string
	files []string
}

var (
	JSON      = jsonRender{}
	XML       = xmlRender{}
	Plain     = plainRender{}
	Redirect  = redirectRender{}
	HTMLDebug = &htmlDebugRender{}
)

func WriteHeader(w http.ResponseWriter, code int, contentType string) {
	w.Header().Set("Content-Type", CombineContentType(contentType))
	w.WriteHeader(code)
}

func CombineContentType(contentType string) string {
	return contentType + "; charset=utf-8"
}

func (_ jsonRender) Render(writer http.ResponseWriter, code int, data ...interface{}) error {
	WriteHeader(writer, code, MIMEJSON)
	encoder := json.NewEncoder(writer)

	return encoder.Encode(data[0])
}

func (_ xmlRender) Render(writer http.ResponseWriter, code int, data ...interface{}) error {
	WriteHeader(writer, code, MIMEXML)
	encoder := xml.NewEncoder(writer)

	return encoder.Encode(data[0])
}

func (_ plainRender) Render(writer http.ResponseWriter, code int, data ...interface{}) error {
	WriteHeader(writer, code, MIMEPlain)
	format := data[0].(string)
	args := data[1].([]interface{})
	var err error
	if len(args) > 0 {
		_, err = writer.Write([]byte(fmt.Sprintf(format, args...)))
	} else {
		_, err = writer.Write([]byte(format))
	}

	return err
}

func (_ redirectRender) Render(writer http.ResponseWriter, code int, data ...interface{}) error {
	writer.WriteHeader(code)
	writer.Header().Set("Location", data[0].(string))

	return nil
}

func (r HTMLRender) Render(writer http.ResponseWriter, code int, data ...interface{}) error {
	WriteHeader(writer, code, MIMEHTML)
	file := data[0].(string)
	obj := data[1]

	return r.Template.ExecuteTemplate(writer, file, obj)
}

func (html *htmlDebugRender) Render(writer http.ResponseWriter, code int, data ...interface{}) error {
	WriteHeader(writer, code, MIMEHTML)
	file := data[0].(string)
	obj := data[1]

	t := template.New("")

	if len(html.files) > 0 {
		if _, err := t.ParseFiles(html.files...); err != nil {
			return err
		}
	}

	for _, glob := range html.globs {
		if _, err := t.ParseGlob(glob); err != nil {
			return err
		}
	}

	return t.ExecuteTemplate(writer, file, obj)
}

func (html *htmlDebugRender) AddGlob(pattern string) {
	html.globs = append(html.globs, pattern)
}

func (html *htmlDebugRender) AddFiles(files ...string) {
	html.files = append(html.files, files...)
}
