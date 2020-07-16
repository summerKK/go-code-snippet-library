package global

import (
	"github.com/summerKK/go-code-snippet-library/blog-service/pkg/logger"
	"github.com/summerKK/go-code-snippet-library/blog-service/pkg/setting"
)

var (
	ServerSetting   *setting.ServerSettingS
	AppSetting      *setting.AppSettingS
	DatabaseSetting *setting.DatabaseSettingS
	Logger          *logger.Logger
)
