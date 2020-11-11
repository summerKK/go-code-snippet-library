package binding

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"net/http"
	"reflect"
	"strings"
)

type Binding interface {
	Bind(req *http.Request, v interface{}) error
}

type jsonBinding struct{}

type xmlBinding struct{}

var (
	JSON = jsonBinding{}
	XML  = xmlBinding{}
)

func (_ jsonBinding) Bind(req *http.Request, v interface{}) error {
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&v); err == nil {
		return Validate(v)
	} else {
		return err
	}
}

func (_ xmlBinding) Bind(req *http.Request, v interface{}) error {
	decoder := xml.NewDecoder(req.Body)
	if err := decoder.Decode(v); err == nil {
		return Validate(v)
	} else {
		return err
	}
}

// 对提交的字段进行校验
func Validate(value interface{}) error {
	var err error
	typ := reflect.TypeOf(value)
	val := reflect.ValueOf(value)

	// 如果是指针类型
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}

	// 对每个字段都进行校验
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i).Interface()
		// 零值
		zero := reflect.Zero(field.Type).Interface()

		// 嵌套结构体
		// 如果field的值是指针.并且不为空
		if field.Type.Kind() == reflect.Struct || (field.Type.Kind() == reflect.Ptr && !reflect.DeepEqual(zero, fieldValue)) {
			err = Validate(fieldValue)
		}

		// 必填认证
		if strings.Index(field.Tag.Get("binding"), "required") > -1 {
			// 字段为空
			if reflect.DeepEqual(zero, fieldValue) {
				name := field.Name
				if j := field.Tag.Get("json"); j != "" {
					name = j
				}

				if f := field.Tag.Get("form"); f != "" {
					name = f
				}

				return errors.New("Required" + name)
			}
		}
	}

	return err
}
