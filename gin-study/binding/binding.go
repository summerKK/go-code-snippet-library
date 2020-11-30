package binding

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

type Binding interface {
	Bind(req *http.Request, v interface{}) error
}

type jsonBinding struct{}

type xmlBinding struct{}

type formBinding struct{}

var (
	JSON = jsonBinding{}
	XML  = xmlBinding{}
	FORM = formBinding{}
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

func (_ formBinding) Bind(req *http.Request, v interface{}) error {
	if err := req.ParseForm(); err != nil {
		return err
	}
	if err := mapForm(v, req.Form); err != nil {
		return err
	}

	return Validate(v)
}

// 把form隐射到ptr上面
func mapForm(ptr interface{}, form map[string][]string) error {
	// 获取映射的结构体
	typ := reflect.TypeOf(ptr).Elem()
	typValue := reflect.ValueOf(ptr).Elem()
	for i := 0; i < typ.NumField(); i++ {
		typeField := typ.Field(i)
		if structFieldName := typeField.Tag.Get("form"); structFieldName != "" {
			structFieldValue := typValue.Field(i)
			if !structFieldValue.CanSet() {
				continue
			}

			// 不存在对应的字段,直接跳过
			formValue, ok := form[structFieldName]
			if !ok {
				continue
			}

			formValueLen := len(formValue)
			// 如果是字段是 slice 类型
			if structFieldValue.Kind() == reflect.Slice && formValueLen > 0 {
				// 获取 slice 的类型
				structValueKind := structFieldValue.Type().Elem().Kind()
				// 创建结构体
				slice := reflect.MakeSlice(structFieldValue.Type(), formValueLen, formValueLen)
				// 把 form 里面的值设置到ptr上面
				for i := 0; i < formValueLen; i++ {
					if err := setWithProperType(structValueKind, formValue[i], slice.Index(i)); err != nil {
						return err
					}
				}
				typValue.Field(i).Set(slice)
			} else {
				if err := setWithProperType(typeField.Type.Kind(), formValue[0], structFieldValue); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func setWithProperType(valueKind reflect.Kind, v string, structField reflect.Value) error {
	switch valueKind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if v == "" {
			v = "0"
		}
		intVal, err := strconv.Atoi(v)
		if err != nil {
			return err
		}
		structField.SetInt(int64(intVal))

	case reflect.Bool:
		if v == "" {
			v = "false"
		}
		boolVal, err := strconv.ParseBool(v)
		if err != nil {
			return err
		}
		structField.SetBool(boolVal)

	case reflect.Float32:
		if v == "" {
			v = "0.0"
		}
		floatVal, err := strconv.ParseFloat(v, 32)
		if err != nil {
			return err
		}
		structField.SetFloat(floatVal)

	case reflect.Float64:
		if v == "" {
			v = "0.0"
		}
		floatVal, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return err
		}
		structField.SetFloat(floatVal)
	case reflect.String:
		structField.SetString(v)
	}

	return nil
}

// 对提交的字段进行校验
func Validate(value interface{}, parents ...string) error {
	var err error
	typ := reflect.TypeOf(value)
	val := reflect.ValueOf(value)

	// 如果是指针类型
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}

	switch typ.Kind() {
	// 结构体
	case reflect.Struct:
		// 对每个字段都进行校验
		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)

			// 过滤字段
			// 未导出的字段直接忽略
			// val.Field(i).CanInterface()可以判断字段是否导出.
			// CanSet()无法判断,因为CanSet的字段需要可寻址的(`addressable`)
			if field.Tag.Get("form") == "-" || !val.Field(i).CanInterface() {
				continue
			}

			fieldValue := val.Field(i).Interface()
			// 零值
			zero := reflect.Zero(field.Type).Interface()

			fieldType := field.Type.Kind()
			// 必填认证
			if strings.Index(field.Tag.Get("binding"), "required") > -1 {
				// 结构体嵌套
				if fieldType == reflect.Struct {
					if reflect.DeepEqual(zero, fieldValue) {
						errMsg := "Required " + field.Name
						if len(parents) > 0 {
							errMsg += strings.Join(parents, ".")
						}

						return errors.New(errMsg)
					}
					// 验证结构体
					err = Validate(fieldValue, append(parents, field.Name)...)
					if err != nil {
						return err
					}

					// 空值
				} else if reflect.DeepEqual(zero, fieldValue) {
					errMsg := "Required " + field.Name
					if len(parents) > 0 {
						errMsg += " on " + strings.Join(parents, ".")
					}
					return errors.New(errMsg)

					// slice
				} else if fieldType == reflect.Slice {
					err = Validate(fieldValue, append(parents, field.Name)...)
					if err != nil {
						return err
					}
				}
			} else {
				// 判断匿名字段
				if fieldType == reflect.Struct {
					err = Validate(fieldValue, append(parents, field.Name)...)
					if err != nil {
						return err
					}
				}
			}
		}

	case reflect.Slice:
		for i := 0; i < val.Len(); i++ {
			fieldValue := val.Index(i).Interface()
			err = Validate(fieldValue, parents...)
			if err != nil {
				return err
			}
		}

	default:
		return nil
	}

	return err
}
