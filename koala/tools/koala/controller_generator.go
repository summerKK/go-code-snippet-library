package main

import (
	"fmt"
	"github.com/emicklei/proto"
	"github.com/summerKK/go-code-snippet-library/koala/logger"
	"os"
	"path"
)

type controllerGenerator struct {
	service *proto.Service
	message []*proto.Message
	rpc     []*proto.RPC
}

func init() {
	controllerGen := &controllerGenerator{}
	_ = genMgr.Register("controller", controllerGen)
}

func (c *controllerGenerator) Run(opt *option) (err error) {
	reader, err := os.Open(opt.Proto3FileName)
	if err != nil {
		logger.Logger.Infof("controller generator open file %s failed:%v", opt.Proto3FileName, err)
		return
	}
	defer reader.Close()

	parser := proto.NewParser(reader)
	definition, _ := parser.Parse()

	proto.Walk(definition,
		proto.WithService(c.handleService),
		proto.WithMessage(c.handleMessage),
		proto.WithRPC(c.handleRpc),
	)

	//logger.Logger.Infof("parse file:%s succ,rpc:%+v", opt.Proto3FileName, c.rpc[0])

	err = c.parseController(opt)
	if err != nil {
		logger.Logger.Infof("controller generator [Run] parse file %s failed:%v", opt.Proto3FileName, err)
	}

	return
}

func (c *controllerGenerator) parseController(opt *option) (err error) {
	filePath := path.Join(opt.Output, "controller", fmt.Sprintf("%s.go", c.service.Name))
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
	_, _ = fmt.Fprintf(writer, "type %s struct {\n}\n\n", c.service.Name)
	for _, rpc := range c.rpc {
		_, _ = fmt.Fprintf(writer, "func (s *%s) %s(ctx context.Context,in %s) (%s,error) {\nreturn}", c.service.Name, rpc.Name, rpc.RequestType, rpc.ReturnsType)
	}
	_, _ = fmt.Fprintln(writer)
	return
}

func (c *controllerGenerator) handleService(s *proto.Service) {
	c.service = s
}

func (c *controllerGenerator) handleMessage(s *proto.Message) {
	c.message = append(c.message, s)
}

func (c *controllerGenerator) handleRpc(s *proto.RPC) {
	c.rpc = append(c.rpc, s)
}
