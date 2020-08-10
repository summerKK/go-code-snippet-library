- 安装etcd
    - go get github.com/coreos/etcd/clientv3@v3.3.18
    - 如果拉取go-systemd模块的时候报错.直接在go.mod里面添加replace
        - replace github.com/coreos/go-systemd => github.com/coreos/go-systemd/v22 v22.0.0

- 常见问题
    - undefined: grpc.SupportPackageIsVersion6 和 undefined: grpc.ClientConnInterface 解决办法
        - 方法1:    
         > 升级grpc到1.27或以上
                                                                                                                      
           **注意：如果升级后出现了其他报错，如 undefined: resolver.BuildOption 或 undefined: resolver.ResolveNowOption，又必须降低grpc版本到1.26或以下时，请使用方法2**

        - 方法2: 
         > 降级grpc版本到v1.26.0
            
           replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
                                                                                                                                                                                                                                                                                                                                                                                                          
         > 降级protoc-gen-go的版本
           
           注意：使用命令 go get -u github.com/golang/protobuf/protoc-gen-go 的效果是安装最新版的protoc-gen-go
           
           降低protoc-gen-go的具体办法，在终端运行如下命令，这里降低到版本 v1.2.0
           
           GIT_TAG="v1.2.0"
           go get -d -u github.com/golang/protobuf/protoc-gen-go
           # 切到v1.2.0版本
           git -C "$(go env GOPATH)"/src/github.com/golang/protobuf checkout $GIT_TAG
           # cd 到 github.com/golang/protobuf/protoc-gen-go
           go install