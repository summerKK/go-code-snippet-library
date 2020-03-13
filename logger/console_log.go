package logger

import "os"

type ConsoleLogger struct {
	level LogLevelType
}

func NewConsoleLogger(level LogLevelType) LogInterface {
	return &ConsoleLogger{level: level}
}

func (c *ConsoleLogger) SetLevel(level LogLevelType) {
	if level > LogFatal || level < LogDebug {
		level = LogDebug
	}
	c.level = level
}

func (c *ConsoleLogger) Debug(format string, args ...interface{}) {
	writeToFile(os.Stdout, c.level, LogDebug, format, args...)
}

func (c *ConsoleLogger) Trace(format string, args ...interface{}) {
	writeToFile(os.Stdout, c.level, LogTrace, format, args...)
}

func (c *ConsoleLogger) Info(format string, args ...interface{}) {
	writeToFile(os.Stdout, c.level, LogInfo, format, args...)
}

func (c *ConsoleLogger) Warn(format string, args ...interface{}) {
	writeToFile(os.Stdout, c.level, LogWarn, format, args...)
}

func (c *ConsoleLogger) Error(format string, args ...interface{}) {
	writeToFile(os.Stdout, c.level, LogError, format, args...)
}

func (c *ConsoleLogger) Fatal(format string, args ...interface{}) {
	writeToFile(os.Stdout, c.level, LogFatal, format, args...)
}

func (c *ConsoleLogger) Close() {

}
