package etcd

import (
	"context"
	"fmt"
	"github.com/summerKK/go-code-snippet-library/koala/registry"
	"testing"
	"time"
)

func TestEtcdRegistry(t *testing.T) {
	// 初始化注册中心
	initRegistry, err := registry.InitRegistry(context.TODO(), "etcd",
		registry.WithTimeout(5*time.Second),
		registry.WithAddrs([]string{"127.0.0.1:2379"}),
		registry.WithHeartBet(5),
		registry.WithRegistryPath("/summer/koala"),
	)

	if err != nil {
		t.Error(err)
		return
	}
	serviceName := "comment_service"
	service := &registry.Service{
		Name: serviceName,
		Nodes: []*registry.Node{
			{
				Id:   0,
				Ip:   "127.0.0.1",
				Port: 10086,
			},
			{
				Id:   1,
				Ip:   "127.0.0.2",
				Port: 10086,
			},
		},
	}
	// 注册服务
	err = initRegistry.Register(context.TODO(), service)
	if err != nil {
		t.Error(err)
	}

	time.Sleep(time.Second)
	ticker := time.After(time.Second * 10)
	for {
		select {
		case <-ticker:
			goto end
		default:
			getService, err := initRegistry.GetService(context.TODO(), serviceName)
			if err != nil {
				t.Fatal(err)
			}
			for i, node := range getService.Nodes {
				fmt.Printf("service:%s node:%d,%+v ", getService.Name, i, node)
			}
			fmt.Println()
			time.Sleep(time.Second)
		}
	}
end:
}
