package etcd

import (
	"context"
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
	service := &registry.Service{
		Name: "comment_service",
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
	err = initRegistry.Register(context.TODO(), service)
	if err != nil {
		t.Error(err)
	}

	time.Sleep(time.Second * 5)
}
