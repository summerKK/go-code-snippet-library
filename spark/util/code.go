package util

const (
	ErrParameters = 1001
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
	default:
		msg = "未知错误"
	}
	return
}
