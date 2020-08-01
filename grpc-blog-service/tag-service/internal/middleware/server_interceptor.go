package middleware

import (
	"context"
	"log"
	"runtime/debug"
	"time"

	"github.com/summerKK/go-code-snippet-library/grpc-blog-service/tag-service/pkg/errcode"
	"google.golang.org/grpc"
)

func ErrLog(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	resp, err = handler(ctx, req)
	if err != nil {
		errLog := "error log: method:%s,code:%v,message:%v,details:%v"
		status := errcode.FromError(err)
		log.Printf(errLog, info.FullMethod, status.Code(), status.Err().Error(), status.Details())
	}

	return resp, err
}

func AccessLog(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	requestLog := "access request log: method:%s begin_time:%d,request:%v"
	beginTime := time.Now().Local().Unix()
	log.Printf(requestLog, info.FullMethod, beginTime, req)

	resp, err = handler(ctx, req)

	responseLog := "access response log: method:%s end_time:%d,response:%v"
	endTime := time.Now().Local().Unix()
	log.Printf(responseLog, info.FullMethod, endTime, resp)

	return resp, err
}

func Recovery(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			recoveryLog := "recovery log: method:%s,message:%v \nstack:%s"
			log.Printf(recoveryLog, info.FullMethod, r, string(debug.Stack()[:]))
		}
	}()

	resp, err = handler(ctx, req)

	return resp, err
}
