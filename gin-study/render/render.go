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

type HTMLRender struct {
	Template *template.Template
}

var (
	JSON  = jsonRender{}
	XML   = xmlRender{}
	Plain = plainRender{}
)

func WriteHeader(w http.ResponseWriter, code int, contentType string) {
	if code >= 0 {
		w.Header().Set("Content-Type", contentType)
		w.WriteHeader(code)
	}
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

func (r HTMLRender) Render(writer http.ResponseWriter, code int, data ...interface{}) error {
	WriteHeader(writer, code, MIMEHTML)
	file := data[0].(string)
	obj := data[1]

	return r.Template.ExecuteTemplate(writer, file, obj)
}

func (_ plainRender) Render(writer http.ResponseWriter, code int, data ...interface{}) error {
	WriteHeader(writer, code, MIMEPlain)
	format := data[0].(string)
	args := data[1].([]interface{})
	_, err := writer.Write([]byte(fmt.Sprintf(format, args...)))

	return err
}
