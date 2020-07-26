package main

import (
	"net"

	pb "github.com/summerKK/go-code-snippet-library/grpc-blog-service/tag-service/proto"
	tagServer "github.com/summerKK/go-code-snippet-library/grpc-blog-service/tag-service/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	server := grpc.NewServer()
	pb.RegisterTagServiceServer(server, tagServer.NewTagServer())
	reflection.Register(server)

	lis, err := net.Listen("tcp", ":8001")
	if err != nil {
		panic(err)
	}

	err = server.Serve(lis)
	if err != nil {
		panic(err)
	}
}
