package main

import (
	"fmt"
	"github.com/summerKK/go-code-snippet-library/koala/logger"
	"os"
	"path"
)

type mainGenerator struct {
}

func init() {
	_ = genMgr.Register("main", &mainGenerator{})
}

func (m *mainGenerator) Run(opt *option) (err error) {
	filePath := path.Join(opt.Output, "main/main.go")
	writer, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		logger.Logger.Infof("main generator [Run] open file %s failed:%v", filePath, err)
	}
	_, _ = fmt.Fprint(writer,
		`
package main

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
`)
	return
}
