package main

import (
	"flag"
	"log"
	"net"
	"net/http"

	"github.com/soheilhy/cmux"
	pb "github.com/summerKK/go-code-snippet-library/grpc-blog-service/tag-service/proto"
	tagServer "github.com/summerKK/go-code-snippet-library/grpc-blog-service/tag-service/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	port string
)

func main() {
	flag.StringVar(&port, "port", "8001", "服务启动端口")
	flag.Parse()

	tcpServer, err := RunTcpServer()
	if err != nil {
		log.Fatalf("Run tcp server error:%v", err)
	}

	mux := cmux.New(tcpServer)
	grpcL := mux.MatchWithWriters(cmux.HTTP2MatchHeaderFieldPrefixSendSettings("content-type", "application/grpc"))
	httpL := mux.Match(cmux.HTTP1Fast())

	grpcServer := RunGrpcServer()
	httpServer := RunHttpServer()
	go grpcServer.Serve(grpcL)
	go httpServer.Serve(httpL)

	err = mux.Serve()
	if err != nil {
		log.Fatalf("mux server run error:%v", err)
	}
}

func RunHttpServer() *http.Server {
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/ping", func(writer http.ResponseWriter, r *http.Request) {
		_, _ = writer.Write([]byte("pong"))
	})

	return &http.Server{
		Addr:    ":" + port,
		Handler: serveMux,
	}
}

func RunGrpcServer() *grpc.Server {
	server := grpc.NewServer()
	pb.RegisterTagServiceServer(server, tagServer.NewTagServer())
	reflection.Register(server)

	return server
}

func RunTcpServer() (net.Listener, error) {
	return net.Listen("tcp", ":"+port)
}
