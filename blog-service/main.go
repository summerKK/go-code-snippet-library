package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/summerKK/go-code-snippet-library/blog-service/global"
	"github.com/summerKK/go-code-snippet-library/blog-service/internal/model"
	"github.com/summerKK/go-code-snippet-library/blog-service/internal/routers"
	"github.com/summerKK/go-code-snippet-library/blog-service/pkg/logger"
	"github.com/summerKK/go-code-snippet-library/blog-service/pkg/setting"
	"github.com/summerKK/go-code-snippet-library/blog-service/pkg/tracer"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	port     string
	runModel string
	config   string

	isVersion    bool
	buildTime    string
	buildVersion string
	gitCommitId  string
)

func init() {
	// 解析参数
	err := setupFlag()
	if err != nil {
		log.Fatalf("init.setupFlag error:%v", err)
	}

	//  读取配置文件
	err = setupSetting()
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

	// 日志监控组件
	err = setupTracer()
	if err != nil {
		log.Fatalf("init.setupTracer error:%v", err)
	}
}

// @title 博客系统
// @version 1.0
// description Go博客
// termsOfService github.com/summerKK/go-code-snippet-library/blog-service
func main() {

	// 打印版本
	if isVersion {
		fmt.Printf("build_time:%s\n", buildTime)
		fmt.Printf("build_version:%s\n", buildVersion)
		fmt.Printf("git_commit_id:%s\n", gitCommitId)
		return
	}

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

	global.Logger.Infof(context.Background(), "%s go-programming-tour-book/%s", "summer", "blog-service")

	go func() {
		err := s.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("s.ListenAndServe error:%v", err)
		}
	}()

	// 等待终端信号
	quit := make(chan os.Signal)
	// 接收syscall.SINGINT和syscall.SIGTREM信号
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down server")

	// 最大时间控制,通知该服务有5s的时间来处理原来的请求
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatal("Server force to shutdown:", err)
	}
	log.Println("server exiting")
}

func setupSetting() error {
	settingS, err := setting.NewSetting(strings.Split(config, ",")...)
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

	err = settingS.ReadSection("JWT", &global.JWTSetting)
	if err != nil {
		return err
	}

	err = settingS.ReadSection("Email", &global.EmailSetting)
	if err != nil {
		return err
	}

	global.ServerSetting.ReadTimeOut *= time.Second
	global.ServerSetting.WriteTimeOut *= time.Second

	if port != "" {
		global.ServerSetting.HttpPort = port
	}

	if runModel != "" {
		global.ServerSetting.RunModel = runModel
	}

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

func setupTracer() error {
	jaegerTracer, _, err := tracer.NewJaegerTracer("blog-service", "127.0.0.1:6831")
	if err != nil {
		return err
	}
	global.Tracer = jaegerTracer

	return nil
}

func setupFlag() error {
	flag.StringVar(&port, "port", "", "启动端口")
	flag.StringVar(&runModel, "model", "", "启动模式")
	flag.StringVar(&config, "config", "configs/", "指定使用的配置文件路径")
	flag.BoolVar(&isVersion, "version", false, "编译信息")
	flag.Parse()

	return nil
}
