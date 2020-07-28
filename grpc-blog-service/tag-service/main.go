package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	pb "github.com/summerKK/go-code-snippet-library/grpc-blog-service/tag-service/proto"
	tagServer "github.com/summerKK/go-code-snippet-library/grpc-blog-service/tag-service/server"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	port string
)

func main() {
	flag.StringVar(&port, "port", "8001", "服务启动端口")
	flag.Parse()

	err := RunTcpServer()
	if err != nil {
		log.Fatalf("RunTcpServer error:%v", err)
	}
}

func RunHttpServer() *http.ServeMux {
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/ping", func(writer http.ResponseWriter, r *http.Request) {
		_, _ = writer.Write([]byte("pong"))
	})

	return serveMux
}

func RunGrpcServer() *grpc.Server {
	server := grpc.NewServer()
	pb.RegisterTagServiceServer(server, tagServer.NewTagServer())
	reflection.Register(server)

	return server
}

func RunTcpServer() error {
	httpMux := RunHttpServer()

	// grpc service 服务提供 restful api支持
	gatewayMux := RunGrpcGatewayServer()
	httpMux.Handle("/", gatewayMux)

	return http.ListenAndServe(":"+port, grpcHandleFunc(RunGrpcServer(), httpMux))
}

func RunGrpcGatewayServer() *runtime.ServeMux {
	endpoint := "0.0.0.0:" + port
	serveMux := runtime.NewServeMux()
	options := []grpc.DialOption{grpc.WithInsecure()}
	_ = pb.RegisterTagServiceHandlerFromEndpoint(context.Background(), serveMux, endpoint, options)

	return serveMux
}

func grpcHandleFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 通过判断是 http1.1(http) http2(grpc) 走不通的handle处理请求
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	}), &http2.Server{})
}
