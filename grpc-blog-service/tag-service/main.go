package main

import (
	"flag"
	"log"
	"net"
	"net/http"

	pb "github.com/summerKK/go-code-snippet-library/grpc-blog-service/tag-service/proto"
	tagServer "github.com/summerKK/go-code-snippet-library/grpc-blog-service/tag-service/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	gRPCPort string
	httpPort string
)

func main() {
	flag.StringVar(&gRPCPort, "gRPC_port", "8001", "gRPC启动端口号")
	flag.StringVar(&httpPort, "http_port", "9001", "http启动端口号")
	flag.Parse()

	errs := make(chan error)
	go func() {
		err := RunHttpServer()
		if err != nil {
			errs <- err
		}
	}()

	go func() {
		err := RunGrpcServer()
		if err != nil {
			errs <- err
		}
	}()

	select {
	case <-errs:
		log.Fatalf("Run server got error:%v", errs)
	}
}

func RunHttpServer() error {
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/ping", func(writer http.ResponseWriter, r *http.Request) {
		_, _ = writer.Write([]byte("pong"))
	})

	return http.ListenAndServe(":"+httpPort, serveMux)
}

func RunGrpcServer() error {
	server := grpc.NewServer()
	pb.RegisterTagServiceServer(server, tagServer.NewTagServer())
	reflection.Register(server)

	lis, err := net.Listen("tcp", ":"+gRPCPort)
	if err != nil {
		return err
	}

	return server.Serve(lis)
}
