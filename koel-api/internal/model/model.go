package model

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/summerKK/go-code-snippet-library/koel-api/global"
	"github.com/summerKK/go-code-snippet-library/koel-api/pkg/setting"
)

type Model struct {
	ID        uint32    `json:"id,omitempty" gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// 初始化数据库
func NewDbEngine(dbSetting *setting.DatabaseSettingS) (*gorm.DB, error) {
	format := "%s:%s@tcp(%s)/%s?charset=%s&parseTime=%t&loc=Local"
	db, err := gorm.Open(dbSetting.DBType, fmt.Sprintf(format,
		dbSetting.Username,
		dbSetting.Password,
		dbSetting.Host,
		dbSetting.DBName,
		dbSetting.Charset,
		dbSetting.ParseTime,
	))

	if err != nil {
		return nil, err
	}

	if global.ServerSetting.RunModel == "debug" {
		db.LogMode(true)
	}

	// 设置连接池最大连接数和空闲数
	db.DB().SetMaxIdleConns(global.DatabaseSetting.MaxIdleConns)
	db.DB().SetMaxOpenConns(global.DatabaseSetting.MaxOpenConns)

	return db, nil
}
