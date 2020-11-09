package binding

import (
	"encoding/json"
	"encoding/xml"
	"io"
)

type Binding interface {
	Bind(reader io.Reader, v interface{}) error
}

type jsonBinding struct{}

type xmlBinding struct{}

var (
	JSON = jsonBinding{}
	XML  = xmlBinding{}
)

func (_ jsonBinding) Bind(reader io.Reader, v interface{}) error {
	return json.NewDecoder(reader).Decode(&v)
}

func (_ xmlBinding) Bind(reader io.Reader, v interface{}) error {
	return xml.NewDecoder(reader).Decode(&v)
}
