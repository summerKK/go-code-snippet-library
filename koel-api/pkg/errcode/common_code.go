package errcode

var (
	Success                   = NewError(0, "成功")
	UnauthorizedAuthNotExist  = NewError(1000003, "鉴权失败,找不到对应的账号或者密码错误")
	UnauthorizedTokenError    = NewError(1000004, "鉴权失败,Token错误")
	UnauthorizedTokenTimeout  = NewError(1000005, "鉴权失败,Token超时")
	UnauthorizedTokenGenerate = NewError(1000006, "鉴权失败,Token生成失败")
	InvalidParams             = NewError(1000001, "入参错误")
)
