# https://segmentfault.com/a/1190000021456180
#protoc --proto_path=../ --go_out=plugins=grpc,paths=source_relative:../ ../proto/*.proto

protoc --proto_path=../ -I/Users/summer/go/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --grpc-gateway_out=logtostderr=true,paths=source_relative:../ ../proto/*.proto
