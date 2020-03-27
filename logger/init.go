package logger

import "errors"

func Init(mod LogModType, level LogLevelType, options ...string) (LogInterface, error) {
	switch mod {
	case FileMod:
		if len(options) < 2 {
			return nil, errors.New("not enought params")
		}
		return NewFileLog(options[0], options[1], level), nil
	case ConsoleMod:
		return NewConsoleLogger(level), nil
	default:
		return nil, errors.New("unsupported mod type")
	}
}
