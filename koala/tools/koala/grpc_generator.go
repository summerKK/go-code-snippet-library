package main

import (
	"fmt"
	"os/exec"
	"path"
)

type grpcGenerator struct {
}

func init() {
	_ = genMgr.Register("gprc", &grpcGenerator{})
}

func (g *grpcGenerator) Run(opt *option, metaData *metaDataService) (err error) {
	joinPath := path.Join(opt.Output, "generate")
	s := fmt.Sprintf("plugins=grpc:%s", joinPath)
	//protoc --go_out=plugins=grpc:. hello.proto
	cmd := exec.Command("protoc", "--go_out", s, opt.Proto3FileName)
	err = cmd.Run()
	return
}
