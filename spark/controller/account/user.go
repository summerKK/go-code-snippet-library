package account

import (
	"github.com/gin-gonic/gin"
	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/summerKK/go-code-snippet-library/spark/common"
	dbaccount "github.com/summerKK/go-code-snippet-library/spark/dal/db/account"
	middlewareAccount "github.com/summerKK/go-code-snippet-library/spark/middleware/account"
	"github.com/summerKK/go-code-snippet-library/spark/service/account"
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
		validation.Field(&userInfo.Sex, validation.In(1, 2)),
	)
	if err != nil {
		util.ResponseErr(c, &util.CodeInfo{
			Code: util.ErrParameters,
			Msg:  err.Error(),
		})
		return
	}

	err = account.Register(&userInfo)
	if err != nil {
		if err == dbaccount.DbUserExists {
			util.ResponseErr(c, &util.CodeInfo{Code: util.ErrUserExists})
		} else {
			util.ResponseErr(c, &util.CodeInfo{Code: util.ErrServiceBusy})
		}
		return
	}

	util.ResponseSuc(c, nil, nil)
}

func Login(c *gin.Context) {
	var err error
	var info *common.UserInfo
	defer func() {
		if err != nil || info == nil {
			return
		}
		// 把用户id存入session
		_ = middlewareAccount.SetUserId(c, int64(info.UserId))
		// 设置cookie
		middlewareAccount.ProcessResponse(c)
		util.ResponseSuc(c, info, nil)
	}()
	// 加载请求中间件
	middlewareAccount.ProcessRequest(c)
	var userInfo common.UserInfo
	err = c.BindJSON(&userInfo)
	if err != nil {
		util.ResponseErr(c, &util.CodeInfo{Code: util.ErrParameters})
		return
	}

	err = validation.ValidateStruct(&userInfo,
		validation.Field(&userInfo.User, validation.Required),
		validation.Field(&userInfo.Password, validation.Required),
	)

	if err != nil {
		util.ResponseErr(c, &util.CodeInfo{
			Code: util.ErrParameters,
			Msg:  err.Error(),
		})
		return
	}

	info, err = account.Login(&userInfo)
	if err != nil {
		if err == dbaccount.DbUserLoginFialed {
			util.ResponseErr(c, &util.CodeInfo{Code: util.ErrUserLoginFailed})
		} else {
			util.ResponseErr(c, &util.CodeInfo{Code: util.ErrServiceBusy})
		}
		return
	}
}
