package logger

import (
	"fmt"
	"os"
)

type FIleLogger struct {
	path     string
	filename string
	file     *os.File
	warnFile *os.File
	level    LogLevelType
}

func NewFileLog(path, filename string, level LogLevelType) LogInterface {
	log := &FIleLogger{
		path:     path,
		filename: filename,
		level:    level,
	}
	log.init()

	return log
}

func (f *FIleLogger) init() {
	fileName := fmt.Sprintf("%s/%s.log", f.path, f.filename)
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		panic(fmt.Sprintf("open file failed,error:%v\n", err))
	}
	f.file = file

	fileName = fmt.Sprintf("%s/%s.log.warn", f.path, f.filename)
	file, err = os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		panic(fmt.Sprintf("open file(warn) failed,error:%v\n", err))
	}
	f.warnFile = file

}

func (f *FIleLogger) SetLevel(level LogLevelType) {
	if level > LogFatal || level < LogDebug {
		level = LogDebug
	}
	f.level = level
}

func (f *FIleLogger) Debug(format string, args ...interface{}) {
	writeToFile(f.file, f.level, LogDebug, format, args...)
}

func (f *FIleLogger) Trace(format string, args ...interface{}) {
	writeToFile(f.file, f.level, LogDebug, format, args...)
}

func (f *FIleLogger) Info(format string, args ...interface{}) {
	writeToFile(f.file, f.level, LogInfo, format, args...)
}

func (f *FIleLogger) Warn(format string, args ...interface{}) {
	writeToFile(f.file, f.level, LogWarn, format, args...)
}

func (f *FIleLogger) Error(format string, args ...interface{}) {
	writeToFile(f.warnFile, f.level, LogError, format, args...)
}

func (f *FIleLogger) Fatal(format string, args ...interface{}) {
	writeToFile(f.warnFile, f.level, LogFatal, format, args...)
}

func (f *FIleLogger) Close() {
	f.file.Close()
	f.warnFile.Close()
}
