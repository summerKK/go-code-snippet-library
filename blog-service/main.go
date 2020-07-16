package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/summerKK/go-code-snippet-library/blog-service/global"
	"github.com/summerKK/go-code-snippet-library/blog-service/internal/model"
	"github.com/summerKK/go-code-snippet-library/blog-service/internal/routers"
	"github.com/summerKK/go-code-snippet-library/blog-service/pkg/logger"
	"github.com/summerKK/go-code-snippet-library/blog-service/pkg/setting"
	"gopkg.in/natefinch/lumberjack.v2"
)

func init() {
	//  读取配置文件
	err := setupSetting()
	if err != nil {
		log.Fatalf("init.setupSetting error:%v", err)
	}

	// 初始化日志
	err = setupLogger()
	if err != nil {
		log.Fatalf("init.setupLogger error:%v", err)
	}

	// 初始化数据库
	err = setupDBEngine()
	if err != nil {
		log.Fatalf("init.setupDBEngine error:%v", err)
	}
}

func main() {
	gin.SetMode(global.ServerSetting.RunModel)
	router := routers.NewRouter()

	s := &http.Server{
		Addr:           ":" + global.ServerSetting.HttpPort,
		Handler:        router,
		TLSConfig:      nil,
		ReadTimeout:    global.ServerSetting.ReadTimeOut,
		WriteTimeout:   global.ServerSetting.WriteTimeOut,
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

	log.Printf("[cofnig] serviceSetting: %+v \n", global.ServerSetting)
	log.Printf("[config] appSetting: %+v \n", global.AppSetting)
	log.Printf("[config] databaseSetting: %+v \n", global.DatabaseSetting)

	return nil
}

func setupDBEngine() error {
	var err error
	global.DBEngine, err = model.NewDBEngine(global.DatabaseSetting)

	return err
}

func setupLogger() error {
	global.Logger = logger.NewLogger(&lumberjack.Logger{
		Filename:  global.AppSetting.LogSavePath + "/" + global.AppSetting.LogFileName + global.AppSetting.LogFileExt,
		MaxSize:   600,
		MaxAge:    10,
		LocalTime: true,
	}, "", log.LstdFlags).WithCaller(2)

	return nil
}
