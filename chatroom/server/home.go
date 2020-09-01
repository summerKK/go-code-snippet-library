package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/summerKK/go-code-snippet-library/chatroom/global"
	"github.com/summerKK/go-code-snippet-library/chatroom/logic"
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

func userListHandleFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	userList := logic.Broadcaster.UserList()
	b, err := json.Marshal(userList)
	if err != nil {
		fmt.Fprint(w, "[]")
	} else {
		fmt.Fprint(w, string(b))
	}
}
