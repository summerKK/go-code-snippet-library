package logger

import "testing"

func TestFileLogger(test *testing.T) {
	log := NewFileLog("./", "test_log_file", LogFatal)
	log.Debug("this is a debug test %v", 1)
	log.Fatal("this is a fatal test %v", 1)
	log.Close()
}

func TestConsoleLogger(test *testing.T) {
	log := NewConsoleLogger(LogDebug)
	log.Debug("this is a debug test %v", 1)
	log.Fatal("this is a fatal test %v", 1)
	log.Close()
}
