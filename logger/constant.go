package logger

type LogLevelType int

const (
	LogDebug LogLevelType = iota
	LogTrace
	LogInfo
	LogWarn
	LogError
	LogFatal
)

func GetLevelText(level LogLevelType) string {
	switch level {
	case LogDebug:
		return "DEBUG"
	case LogTrace:
		return "TRACE"
	case LogInfo:
		return "INFO"
	case LogWarn:
		return "WARNING"
	case LogError:
		return "ERROR"
	case LogFatal:
		return "FATAL"
	default:
		return ""
	}
}
