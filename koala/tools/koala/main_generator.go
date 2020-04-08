package main

import (
	"github.com/summerKK/go-code-snippet-library/koala/logger"
	"html/template"
	"os"
	"path"
)

type mainGenerator struct {
}

func init() {
	_ = genMgr.Register("main", &mainGenerator{})
}

func (m *mainGenerator) Run(opt *option, metaData *metaDataService) (err error) {
	filePath := path.Join(opt.Output, "main/main.go")
	writer, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		logger.Logger.Infof("main generator [Run] open file %s failed:%v", filePath, err)
	}

	parse, err := template.New("main").Parse(main_template)
	if err != nil {
		logger.Logger.Infof("main generator [Run] parse template file failed:%v", err)
		return
	}
	err = parse.Execute(writer, metaData)
	if err != nil {
		logger.Logger.Infof("main generator [Run] parse Execute failed:%v", err)
		return
	}

	return
}
