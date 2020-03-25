package account

import "errors"

var (
	DbUserExists      = errors.New("用户已存在")
	DbUserLoginFialed = errors.New("账户或者密码错误")
)
