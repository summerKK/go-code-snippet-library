# https://segmentfault.com/a/1190000021456180
# 生成 *.pb.go 文件
protoc -I../ -I/Users/summer/go/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --go_out=plugins=grpc,paths=source_relative:../ ../proto/*.proto

# 生成 *.pb.gw.go 文件
protoc --proto_path=../ -I/Users/summer/go/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --grpc-gateway_out=logtostderr=true,paths=source_relative:../ ../proto/*.proto

# 生成 swagger.json 文件
protoc --proto_path=../ -I/Users/summer/go/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --swagger_out=logtostderr=true:../ ../proto/*.proto

# 生成 swagger 文档 (命令只能在根目录执行,要不然路径会有问题.)
#go-bindata --nocompress -pkg swagger -o pkg/swagger/data.go third_party/swagger-ui/...
