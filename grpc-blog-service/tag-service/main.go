package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/proxy/grpcproxy"
	assetfs "github.com/elazarl/go-bindata-assetfs"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/summerKK/go-code-snippet-library/grpc-blog-service/tag-service/global"
	"github.com/summerKK/go-code-snippet-library/grpc-blog-service/tag-service/internal/middleware"
	"github.com/summerKK/go-code-snippet-library/grpc-blog-service/tag-service/pkg/swagger"
	"github.com/summerKK/go-code-snippet-library/grpc-blog-service/tag-service/pkg/tracer"
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

const SERVICE_NAME = "rpc-blog-service"

func init() {
	err := setUpTracer()
	if err != nil {
		log.Printf("setUpTracer error:%v", err)
	}
}

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

	prefix := "/swagger-ui/"
	fileServer := http.FileServer(&assetfs.AssetFS{
		Asset:    swagger.Asset,
		AssetDir: swagger.AssetDir,
		Prefix:   "third_party/swagger-ui",
	})
	serveMux.Handle(prefix, http.StripPrefix(prefix, fileServer))
	serveMux.HandleFunc("/swagger/", func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "swagger.json") {
			http.NotFound(w, r)
			return
		}

		p := strings.TrimPrefix(r.URL.Path, "/swagger/")
		p = path.Join("proto", p)

		http.ServeFile(w, r, p)
	})

	return serveMux
}

func RunGrpcServer() *grpc.Server {
	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			middleware.Recovery,
			middleware.AccessLog,
			middleware.ErrLog,
			// 联络追踪
			middleware.ServerTracing,
		)),
	}
	server := grpc.NewServer(opts...)
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

	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: time.Second * 60,
	})
	if err != nil {
		return err
	}
	defer etcdClient.Close()

	target := fmt.Sprintf("/etcdv3://goproject/grpc/%s", SERVICE_NAME)
	grpcproxy.Register(etcdClient, target, ":"+port, 60)

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
		// 捕捉到自定义错误,直接把code和message改为自定义的
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

func setUpTracer() error {
	jaegerTracer, _, err := tracer.NewJaegerTracer(SERVICE_NAME, "127.0.0.1:6831")
	if err != nil {
		return err
	}
	global.Tracer = jaegerTracer

	return nil
}
