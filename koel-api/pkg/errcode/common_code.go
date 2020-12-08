package errcode

var (
	Success                   = NewError(200, "成功")
	UnauthorizedAuthNotExist  = NewError(401, "鉴权失败,找不到对应的账号或者密码错误")
	UnauthorizedTokenError    = NewError(401, "鉴权失败,Token错误")
	UnauthorizedTokenTimeout  = NewError(401, "鉴权失败,Token超时")
	UnauthorizedTokenGenerate = NewError(401, "鉴权失败,Token生成失败")
	InvalidParams             = NewError(400, "入参错误")
)
