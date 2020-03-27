package logger

import "github.com/summerKK/go-code-snippet-library/logger"

var Logger logger.LogInterface

func Init() (err error) {
	Logger, err = logger.Init(logger.ConsoleMod, logger.LogDebug)
	return
}
