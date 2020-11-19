package gin

import "os"

const GINMode = "GIN_MODE"

const (
	DebugMode   string = "debug"
	ReleaseMode string = "release"
)

const (
	debugCode = iota
	releaseCode
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
	default:
		panic("gin mode unknown, the allowed modes are: " + DebugMode + " and " + ReleaseMode)
	}
}
