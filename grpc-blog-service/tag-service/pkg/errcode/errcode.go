package errcode

import (
	"fmt"
	"net/http"
)

type Error struct {
	code int
	msg  string
}

var _codes = map[int]string{}

func NewError(code int, msg string) *Error {
	if _, ok := _codes[code]; ok {
		panic(fmt.Sprintf("code码:%d已存在,请更换一个", code))
	}
	_codes[code] = msg

	return &Error{
		code: code,
		msg:  msg,
	}
}

func (e *Error) Error() string {
	return fmt.Sprintf("错误码:%d,错误信息:%s", e.code, e.msg)
}

func (e *Error) Code() int {
	return e.code
}

func (e *Error) Msg() string {
	return e.msg
}

func (e *Error) StatusCode() int {
	switch e.Code() {
	case Success.Code():
		return http.StatusOK
	default:
		return http.StatusInternalServerError
	}
}
