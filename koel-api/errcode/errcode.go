package errcode

import (
	"fmt"
	"net/http"
)

type Error struct {
	code    int
	msg     string
	details []string
}

var codes = map[int]string{}

func NewError(code int, msg string) *Error {
	if _, ok := codes[code]; ok {
		panic(fmt.Sprintf("错误码 %d 已存在,请更换一个", code))
	}

	codes[code] = msg

	return &Error{
		code: code,
		msg:  msg,
	}
}

func (e *Error) String() string {
	return fmt.Sprintf("错误码:%d,错误信息:%s", e.Code(), e.Msg())
}

func (e *Error) Code() int {
	return e.code
}

func (e *Error) Msg() string {
	return e.msg
}

func (e Error) Details() []string {
	return e.details
}

func (e *Error) Msgf(args []interface{}) string {
	return fmt.Sprintf(e.msg, args...)
}

// 添加额外信息
func (e *Error) WithDetails(details ...string) *Error {
	e.details = make([]string, 0, len(details))
	for _, detail := range details {
		e.details = append(e.details, detail)
	}

	return e
}

// 判断返回服务器的状态码
func (e *Error) StatusCode() int {
	switch e.Code() {
	case Success.Code():
		return http.StatusOK
	case UnauthorizedAuthNotExist.Code():
		fallthrough
	case UnauthorizedTokenError.Code():
		fallthrough
	case UnauthorizedTokenGenerate.Code():
		fallthrough
	case UnauthorizedTokenTimeout.Code():
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}
