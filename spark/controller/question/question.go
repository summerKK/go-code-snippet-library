package question

import (
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/summerKK/go-code-snippet-library/spark/common"
	"github.com/summerKK/go-code-snippet-library/spark/filter"
	"github.com/summerKK/go-code-snippet-library/spark/service/question"
	"github.com/summerKK/go-code-snippet-library/spark/util"
)

func Save(c *gin.Context) {
	var q common.Question
	err := c.BindJSON(&q)
	if err != nil {
		util.ResponseErr(c, &util.CodeInfo{Code: util.ErrParameters})
		return
	}
	err = validation.ValidateStruct(&q,
		validation.Field(&q.CategoryId, validation.Required),
		validation.Field(&q.Content, validation.Required),
		validation.Field(&q.Caption, validation.Required),
	)
	if err != nil {
		util.ResponseErr(c, &util.CodeInfo{
			Code: util.ErrParameters,
			Msg:  err.Error(),
		})
		return
	}

	// 检查是否含有敏感词
	_, replace := filter.Filter(q.Caption, "***")
	if replace {
		util.ResponseErr(c, &util.CodeInfo{
			Code: util.ErrSensitiveWord,
		})
		return
	}

	_, replace = filter.Filter(q.Content, "***")
	if replace {
		util.ResponseErr(c, &util.CodeInfo{
			Code: util.ErrSensitiveWord,
		})
		return
	}

	err = question.Save(c, &q)
	if err != nil {
		util.ResponseErr(c, &util.CodeInfo{Code: util.ErrServiceBusy})
		return
	}

	util.ResponseSuc(c, nil, nil)
}
