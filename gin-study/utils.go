package gin

import (
	"path"
	"reflect"
	"runtime"
)

func joinGroupPath(elems ...string) string {
	joined := path.Join(elems...)
	lastComponent := elems[len(elems)-1]
	if len(lastComponent) > 0 && lastComponent[len(lastComponent)-1] == '/' && joined[len(joined)-1] != '/' {
		joined += "/"
	}

	return joined
}

func filterFlags(content string) string {
	for i, c := range content {
		if c == ' ' || c == ';' {
			return content[:i]
		}
	}

	return content
}

func funcName(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}
