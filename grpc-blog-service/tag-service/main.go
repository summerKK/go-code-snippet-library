package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	pb "github.com/summerKK/go-code-snippet-library/grpc-blog-service/tag-service/proto"
	grpcServer "github.com/summerKK/go-code-snippet-library/grpc-blog-service/tag-service/server"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type HttpError struct {
	Code    int32  `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

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
	// 注册tagService
	pb.RegisterTagServiceServer(server, grpcServer.NewTagServer())
	// 注册articleService
	pb.RegisterArticleServiceServer(server, grpcServer.NewArticleServer())
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
	runtime.HTTPError = grpcGatewayError
	endpoint := "0.0.0.0:" + port
	serveMux := runtime.NewServeMux()
	options := []grpc.DialOption{grpc.WithInsecure()}
	_ = pb.RegisterTagServiceHandlerFromEndpoint(context.Background(), serveMux, endpoint, options)
	_ = pb.RegisterArticleServiceHandlerFromEndpoint(context.Background(), serveMux, endpoint, options)

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

func grpcGatewayError(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
	s, ok := status.FromError(err)
	if !ok {
		s = status.New(codes.Unknown, err.Error())
	}
	httpError := HttpError{
		Code:    int32(s.Code()),
		Message: s.Message(),
	}
	details := s.Details()
	for _, detail := range details {
		if v, ok := detail.(*pb.Error); ok {
			httpError.Code = v.Code
			httpError.Message = v.Message
		}
	}

	marshal, _ := json.Marshal(httpError)
	w.Header().Set("Content-Type", marshaler.ContentType())
	w.WriteHeader(runtime.HTTPStatusFromCode(s.Code()))

	_, _ = w.Write(marshal)
}
