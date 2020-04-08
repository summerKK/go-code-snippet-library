package main

var main_template = `
package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"github.com/summerKK/go-code-snippet-library/koala/tools/koala/output/controller"
	hello "github.com/summerKK/go-code-snippet-library/koala/tools/koala/output/generate/example"
)

const (
	port = ":50051"
)

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	server := controller.Hello{}
	hello.Register{{.Service.Name}}Server(s, &server)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
`
