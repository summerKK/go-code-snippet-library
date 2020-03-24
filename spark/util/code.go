package util

const (
	ErrParameters  = 1001
	ErrUserExists  = 1002
	ErrServiceBusy = 1003
)

const (
	SucDefault = 0
)

func GetCodeMsg(code int) (msg string) {
	switch code {
	case SucDefault:
		msg = "success"
	case ErrParameters:
		msg = "参数错误"
	case ErrUserExists:
		msg = "用户已经存在"
	case ErrServiceBusy:
		msg = "服务器繁忙"
	default:
		msg = "未知错误"
	}
	return
}
