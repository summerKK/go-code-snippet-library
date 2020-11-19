package gin

import "os"

const GINMode = "GIN_MODE"

const (
	DebugMode   string = "debug"
	ReleaseMode string = "release"
	TestModel   string = "test"
)

const (
	debugCode = iota
	releaseCode
	testMode
)

var ginMode int = debugCode

func init() {
	value := os.Getenv(GINMode)
	if value == "" {
		SetMode(DebugMode)
	} else {
		SetMode(value)
	}
}

func SetMode(v string) {
	switch v {
	case DebugMode:
		ginMode = debugCode
	case ReleaseMode:
		ginMode = releaseCode
	case TestModel:
		ginMode = testMode
	default:
		panic("gin mode unknown, the allowed modes are: " + DebugMode + " and " + ReleaseMode)
	}
}
