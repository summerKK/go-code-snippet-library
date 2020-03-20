package account

import (
	"github.com/gin-gonic/gin"
	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/summerKK/go-code-snippet-library/spark/common"
	"github.com/summerKK/go-code-snippet-library/spark/util"
)

func Register(c *gin.Context) {
	var userInfo common.UserInfo
	err := c.BindJSON(&userInfo)
	if err != nil {
		util.ResponseErr(c, &util.CodeInfo{Code: util.ErrParameters})
		return
	}
	err = validation.ValidateStruct(&userInfo,
		validation.Field(&userInfo.User, validation.Required),
		validation.Field(&userInfo.Email, validation.Required, is.Email),
		validation.Field(&userInfo.Nickname, validation.Required),
		validation.Field(&userInfo.Password, validation.Required),
		validation.Field(&userInfo.Sex, validation.Required, validation.In([]interface{}{1, 2})),
	)
	if err != nil {
		util.ResponseErr(c, &util.CodeInfo{
			Code: util.ErrParameters,
			Msg:  err.Error(),
		})
		return
	}

	util.ResponseSuc(c, &userInfo, nil)
}

func Login(c *gin.Context) {

}
