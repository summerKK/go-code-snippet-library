package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	pb "github.com/summerKK/go-code-snippet-library/grpc-demo/proto"
	"google.golang.org/grpc"
)

type GreeterServer struct {
}

func (sv GreeterServer) SayRoute(stream pb.Greeter_SayRouteServer) error {
	for {
		time.Sleep(time.Second)
		err := stream.Send(&pb.HelloReplay{Message: "hello,world"})
		if err != nil {
			return err
		}
		response, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		log.Printf("sayRoute response:%v", response)
	}
}

func (sv GreeterServer) SayRecord(stream pb.Greeter_SayRecordServer) error {
	list := map[int32]*pb.HelloReplay{}
	var i int32 = 0
	for {
		request, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return stream.SendAndClose(&pb.HelloRecord{List: list})
			}
			return err
		}
		list[i] = &pb.HelloReplay{
			Message: request.Name,
		}
		i++
	}
}

func (sv GreeterServer) SayList(request *pb.HelloRequest, server pb.Greeter_SayListServer) error {
	for i := 0; i < 100; i++ {
		time.Sleep(time.Second)
		err := server.Send(&pb.HelloReplay{
			Message: "hello,world " + fmt.Sprintf("%d", i),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (sv GreeterServer) SayHello(ctx context.Context, r *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{
		Message: r.Name + " hello,world",
	}, nil
}

func main() {
	server := grpc.NewServer()
	pb.RegisterGreeterServer(server, GreeterServer{})
	listen, err := net.Listen("tcp", ":8082")
	if err != nil {
		panic(err)
	}

	_ = server.Serve(listen)
}
