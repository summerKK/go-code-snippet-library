package app

import (
	"github.com/gin-gonic/gin"
	"github.com/summerKK/go-code-snippet-library/koel-api/pkg/errcode"
)

type Response struct {
	Ctx *gin.Context
}

func NewResponse(ctx *gin.Context) *Response {
	return &Response{Ctx: ctx}
}

func (r *Response) ToResponse(data interface{}, err *errcode.Error) {
	r.Ctx.JSON(err.StatusCode(), gin.H{
		"code":    err.Code(),
		"message": err.Msg(),
		"data":    data,
		"details": err.Details(),
	})
}

// 列表返回
func (r *Response) Success(data interface{}) {
	r.ToResponse(data, errcode.Success)
}

func (r *Response) ToErrorResponse(err *errcode.Error) {
	r.ToResponse(nil, err)
}

type Pagination struct {
	Page      int `json:"page"`
	PageSize  int `json:"page_size"`
	TotalRows int `json:"total_rows"`
}
