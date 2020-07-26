# https://segmentfault.com/a/1190000021456180
protoc --proto_path=../ --go_out=plugins=grpc,paths=source_relative:../ ../proto/*.proto
