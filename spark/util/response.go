package util

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type response struct {
	Code int                    `json:"code"`
	Msg  string                 `json:"msg"`
	Data map[string]interface{} `json:"data"`
}

type CodeInfo struct {
	Code int
	Msg  string
}

func ResponseErr(c *gin.Context, code *CodeInfo) {
	c.JSON(http.StatusOK, &response{
		Code: code.Code,
		Msg:  codeMsg(code),
		Data: make(map[string]interface{}),
	})
}

func ResponseSuc(c *gin.Context, data interface{}, code *CodeInfo) {
	if code == nil {
		code = &CodeInfo{
			Code: SucDefault,
		}
	}
	resp := &response{
		Code: code.Code,
		Msg:  codeMsg(code),
		Data: make(map[string]interface{}),
	}
	resp.Data["data"] = data
	c.JSON(http.StatusOK, resp)
}

func codeMsg(code *CodeInfo) (msg string) {
	return If(code.Msg != "", code.Msg, GetCodeMsg(code.Code)).(string)
}
