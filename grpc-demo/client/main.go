package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	pb "github.com/summerKK/go-code-snippet-library/grpc-demo/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func main() {
	ctx := context.Background()
	md := metadata.New(map[string]string{"go": "hello,world"})
	ctx = metadata.NewOutgoingContext(ctx, md)
	conn, err := grpc.DialContext(ctx, ":8082", grpc.WithInsecure())
	defer conn.Close()

	if err != nil {
		panic(err)
	}

	//
	client := pb.NewGreeterClient(conn)
	err = sayHello(client)
	if err != nil {
		panic(err)
	}

	// 服务端流式rpc
	//err = sayList(client)
	//if err != nil {
	//	panic(err)
	//}

	// 客户端流式rpc
	//err = sayRecord(client)
	//if err != nil {
	//	panic(err)
	//}

	// 双向流式rpc
	//err = sayRoute(client)
	//if err != nil {
	//	panic(err)
	//}
}

func sayHello(client pb.GreeterClient) error {
	var header metadata.MD
	response, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: "summer"}, grpc.Header(&header))
	if err != nil {
		return err
	}
	log.Println(response)

	log.Printf("receive metadata from context:%v", header)

	return nil
}

func sayList(client pb.GreeterClient) error {
	stream, err := client.SayList(context.Background(), &pb.HelloRequest{Name: "summer"})
	if err != nil {
		return err
	}
	for {
		resp, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		log.Printf("resp: %v", resp)
	}

	return nil
}

func sayRecord(client pb.GreeterClient) error {
	stream, err := client.SayRecord(context.Background())
	if err != nil {
		return err
	}

	for i := 0; i < 100; i++ {
		err := stream.Send(&pb.HelloRequest{
			Name: "summer" + fmt.Sprintf("%d", i),
		})
		if err != nil {
			return err
		}
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		return err
	}

	log.Printf("say Record response:%v", resp)

	return nil
}

func sayRoute(client pb.GreeterClient) error {
	stream, err := client.SayRoute(context.Background())
	if err != nil {
		return err
	}

	for i := 0; i < 100; i++ {
		time.Sleep(time.Second)
		err := stream.Send(&pb.HelloRequest{
			Name: "summer" + fmt.Sprintf("%d", i),
		})
		if err != nil {
			return err
		}

		response, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		log.Printf("sayRoute response:%v", response)
	}

	err = stream.CloseSend()
	if err != nil {
		return err
	}

	return nil
}
