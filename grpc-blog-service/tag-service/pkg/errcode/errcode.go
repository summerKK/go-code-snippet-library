package errcode

import (
	"fmt"
	"net/http"
)

type Error struct {
	code    int      `json:"code"`
	msg     string   `json:"msg"`
	details []string `json:"details"`
}

var codes map[int]string

func NewError(code int, msg string) *Error {
	if _, ok := codes[code]; ok {
		panic(fmt.Sprintf("code码:%d已存在,请更换一个", code))
	}
	codes[code] = msg

	return &Error{
		code: code,
		msg:  msg,
	}
}

func (e *Error) String() string {
	return fmt.Sprintf("错误码:%d,错误信息:%s", e.code, e.msg)
}

func (e *Error) Code() int {
	return e.code
}

func (e *Error) Msg() string {
	return e.msg
}

func (e *Error) Msgf(args []interface{}) string {
	return fmt.Sprintf(e.msg, args...)
}

func (e *Error) Details() []string {
	return e.details
}

func (e *Error) WithDetails(details ...string) *Error {
	e.details = make([]string, 0, len(details))
	for _, detail := range details {
		e.details = append(e.details, detail)
	}

	return e
}

func (e *Error) StatusCode() int {
	switch e.Code() {
	case Success.Code():
		return http.StatusOK
	default:
		return http.StatusInternalServerError
	}
}
