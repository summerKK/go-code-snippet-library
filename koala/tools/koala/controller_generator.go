package main

import (
	"fmt"
	"github.com/summerKK/go-code-snippet-library/koala/logger"
	"os"
	"path"
)

type controllerGenerator struct {
}

func init() {
	controllerGen := &controllerGenerator{}
	_ = genMgr.Register("controller", controllerGen)
}

func (c *controllerGenerator) Run(opt *option, metaData *metaDataService) (err error) {
	//logger.Logger.Infof("parse file:%s succ,Rpc:%+v", opt.Proto3FileName, c.Rpc[0])

	err = c.parseController(opt, metaData)
	if err != nil {
		logger.Logger.Infof("controller generator [Run] parse file %s failed:%v", opt.Proto3FileName, err)
	}

	return
}

func (c *controllerGenerator) parseController(opt *option, metaData *metaDataService) (err error) {
	filePath := path.Join(opt.Output, "controller", fmt.Sprintf("%s.go", metaData.Service.Name))
	writer, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		logger.Logger.Infof("controller generator [parseController] open file %s failed:%v", filePath, err)
		return
	}
	/*
		package main

		import (
			"context"
			"log"
			"net"

			"google.golang.org/grpc"
			pb "google.golang.org/grpc/examples/helloworld/helloworld"
		)

		// server is used to implement helloworld.GreeterServer.
		type server struct {
			pb.UnimplementedGreeterServer
		}

		// SayHello implements helloworld.GreeterServer
		func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
			log.Printf("Received: %v", in.GetName())
			return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
		}
	*/
	_, _ = fmt.Fprintf(writer, "package controller\n")
	_, _ = fmt.Fprint(writer, "import (\n")
	_, _ = fmt.Fprint(writer, "\"context\"\n)\n\n")
	_, _ = fmt.Fprintf(writer, "type %s struct {\n}\n\n", metaData.Service.Name)
	for _, rpc := range metaData.Rpc {
		_, _ = fmt.Fprintf(writer, "func (s *%s) %s(ctx context.Context,in %s) (%s,error) {\nreturn}", metaData.Service.Name, rpc.Name, rpc.RequestType, rpc.ReturnsType)
	}
	_, _ = fmt.Fprintln(writer)
	return
}
