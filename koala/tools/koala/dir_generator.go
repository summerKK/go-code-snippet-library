package main

import (
	"github.com/summerKK/go-code-snippet-library/koala/logger"
	"os"
	"path"
)

var dirList = []string{
	"controller",
	"idl",
	"app",
	"conf",
	"generate",
	"main",
	"app/router",
	"app/config",
	"scripts",
}

type dirGenerator struct {
	dirList []string
}

func init() {
	_ = genMgr.Register("dir", &dirGenerator{dirList: dirList})
}

func (d *dirGenerator) Run(opt *option) (err error) {
	for _, dir := range d.dirList {
		joinPath := path.Join(opt.Output, dir)
		err = os.MkdirAll(joinPath, 0755)
		if err != nil {
			logger.Logger.Infof("dir generator make file failed:%v", err)
			return
		}
	}
	return
}
