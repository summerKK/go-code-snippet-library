package util

const (
	ErrParameters = 1000 + iota
	ErrUserExists
	ErrUserLoginFailed
	ErrServiceBusy
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
	case ErrUserLoginFailed:
		msg = "用户名或者密码错误"
	default:
		msg = "未知错误"
	}
	return
}
