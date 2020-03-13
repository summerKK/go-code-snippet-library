package logger

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"time"
)

func GetCallStack() (file, funName string, line int) {
	pc, file, line, ok := runtime.Caller(3)
	if ok {
		funName = runtime.FuncForPC(pc).Name()
	}
	return
}

func writeToFile(file *os.File, currentLevel LogLevelType, level LogLevelType, format string, args ...interface{}) {
	if currentLevel > level {
		return
	}

	timeStr := time.Now().Format("2006-01-02 15:04:05.999")
	errorStr := fmt.Sprintf(format, args...)
	fileName, funName, line := GetCallStack()

	fmt.Fprintf(file, "%s %s (%s %s:%d) %s\n", timeStr, GetLevelText(level), path.Base(fileName), path.Base(funName), line, errorStr)
}
