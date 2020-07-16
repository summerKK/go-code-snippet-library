package main

import (
	"log"
	"net/http"
	"time"

	"github.com/summerKK/go-code-snippet-library/blog-service/global"
	"github.com/summerKK/go-code-snippet-library/blog-service/internal/routers"
	"github.com/summerKK/go-code-snippet-library/blog-service/pkg/setting"
)

func init() {
	err := setupSetting()
	if err != nil {
		log.Fatalf("init.setupSetting error:%v", err)
	}
}

func main() {
	router := routers.NewRouter()

	s := &http.Server{
		Addr:           ":8008",
		Handler:        router,
		TLSConfig:      nil,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	err := s.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func setupSetting() error {
	settingS, err := setting.NewSetting()
	if err != nil {
		return err
	}
	err = settingS.ReadSection("Server", &global.ServerSetting)
	if err != nil {
		return err
	}

	err = settingS.ReadSection("App", &global.AppSetting)
	if err != nil {
		return err
	}

	err = settingS.ReadSection("Database", &global.DatabaseSetting)
	if err != nil {
		return err
	}

	global.ServerSetting.ReadTimeOut *= time.Second
	global.ServerSetting.WriteTimeOut *= time.Second

	return nil
}
