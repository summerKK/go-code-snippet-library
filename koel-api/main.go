package main

import (
	"log"

	"github.com/summerKK/go-code-snippet-library/koel-api/global"
	"github.com/summerKK/go-code-snippet-library/koel-api/internal/model"
	"github.com/summerKK/go-code-snippet-library/koel-api/pkg/setting"
)

var (
	configPath = []string{"configs/"}
)

func init() {
	var err error
	// 读取配置文件
	err = setupSetting()
	if err != nil {
		log.Fatalf("init.Setting error:%v", err)
	}

	// 初始化数据库
	err = setupDBEngine()
	if err != nil {
		log.Fatalf("init.DBEngine error:%v", err)
	}
}

func main() {

}

func setupSetting() error {
	var err error
	SettingS, err := setting.NewSetting(configPath...)
	if err != nil {
		return err
	}

	err = SettingS.ReadSection("Database", &global.DatabaseSetting)
	if err != nil {
		return err
	}

	err = SettingS.ReadSection("Server", &global.ServerSetting)
	if err != nil {
		return err
	}

	return nil
}

func setupDBEngine() error {
	var err error
	global.DBEngine, err = model.NewDbEngine(global.DatabaseSetting)

	return err
}
