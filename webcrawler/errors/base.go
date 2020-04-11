package errors

type ErrorType string

const (
	ERROR_TYPE_DOWNLOADER ErrorType = "downloader error"
	ERROR_TYPE_ANALYZER   ErrorType = "analyzer error"
	ERROR_TYPE_PIPELINE   ErrorType = "pipeline error"
	ERROR_TYPE_SCHEDULER  ErrorType = "scheduler error"
)

type ICrawlerError interface {
	// 错误类型
	Type() ErrorType
	Error() string
}
