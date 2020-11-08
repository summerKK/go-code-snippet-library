package gin

import (
	"errors"
	"reflect"
	"strings"
)

// 对提交的字段进行校验
func Validate(c *Context, value interface{}) error {
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
			err = Validate(c, fieldValue)
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

				err = errors.New("Required" + name)
				c.Error(err, "json validation")
			}
		}
	}

	return err
}
