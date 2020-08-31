package server

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/summerKK/go-code-snippet-library/chatroom/global"
)

func homeHandleFunc(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.ParseFiles(global.RootDir + "/template/home.html")
	if err != nil {
		fmt.Fprint(w, "模板解析错误")
		return
	}

	err = tpl.Execute(w, nil)
	if err != nil {
		fmt.Fprint(w, "模板执行错误")
		return
	}
}
