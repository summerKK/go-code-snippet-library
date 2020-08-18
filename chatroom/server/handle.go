package server

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/summerKK/go-code-snippet-library/chatroom/logic"
)

var (
	rootDir string
)

func RegisterHandle() {
	inferRootDir()

	// 广播消息处理
	go logic.Broadcast.Start()

	http.HandleFunc("/", homeHandleFunc)
	http.HandleFunc("/ws", websocketHandleFunc)
}

// 推断出项目根目录
func inferRootDir() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	var infer func(d string) string
	infer = func(d string) string {
		if exist(d + "/template") {
			return d
		}

		return infer(filepath.Dir(d))
	}

	rootDir = infer(cwd)
}

func exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}
