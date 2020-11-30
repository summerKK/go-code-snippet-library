package gin

import (
	"fmt"
	"path"
	"reflect"
	"runtime"
	"strings"
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

func nameOfFunction(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

func chooseData(custom, wildcard interface{}) interface{} {
	if custom == nil {
		if wildcard == nil {
			panic("negotiation config is invalid")
		}
		return wildcard
	}
	return custom
}

func parseAccept(accept string) []string {
	parts := strings.Split(accept, ",")

	for i, part := range parts {
		index := strings.IndexByte(part, ';')
		if index >= 0 {
			part = accept[0:index]
		}
		part = strings.TrimSpace(part)
		parts[i] = part
	}

	return parts
}

func debugPrint(format string, values ...interface{}) {
	if IsDebugging() {
		fmt.Printf("[GIN-debug] "+format, values...)
	}
}

func lastChar(str string) uint8 {
	size := len(str)
	if size == 0 {
		panic("The length of the string can't be 0")
	}

	return str[size-1]
}
