package errcode

var (
	Success                   = NewError(0, "成功")
	ServerError               = NewError(1000000, "服务器内部错误")
	InvalidParams             = NewError(1000001, "入参错误")
	NotFound                  = NewError(1000002, "找不到")
	UnauthorizedAuthNotExist  = NewError(1000003, "鉴权失败,找不到对应的Appkey和AppSecret")
	UnauthorizedTokenError    = NewError(1000004, "鉴权失败,Token错误")
	UnauthorizedTokenTimeout  = NewError(1000005, "鉴权失败,Token超时")
	UnauthorizedTokenGenerate = NewError(1000006, "鉴权失败,Token生成失败")
	TooManyRequests           = NewError(1000007, "请求过多")
)
