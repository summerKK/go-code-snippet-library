package etcd

import (
	"context"
	"fmt"
	"github.com/summerKK/go-code-snippet-library/koala/registry"
	"testing"
	"time"
)

var (
	initRegistry registry.IRegistry
	serviceName0 = "comment_service_0"
)

func init() {
	var err error
	// 初始化注册中心
	initRegistry, err = registry.InitRegistry(context.TODO(), "etcd",
		registry.WithTimeout(5*time.Second),
		registry.WithAddrs([]string{"127.0.0.1:2379"}),
		registry.WithHeartBet(5),
		registry.WithRegistryPath("/summer/koala"),
	)

	if err != nil {
		panic(err)
	}
	service0 := &registry.Service{
		Name: serviceName0,
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
	err = initRegistry.Register(context.TODO(), service0)
	if err != nil {
		panic(err)
	}
}

func TestEtcdRegistry(t *testing.T) {

	var err error
	// 在 `comment_service` 增加一个新节点
	go func() {
		time.Sleep(time.Second * 5)

		serviceName1 := "comment_service_0"
		service1 := &registry.Service{
			Name: serviceName1,
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
				{
					Id:   2,
					Ip:   "127.0.0.3",
					Port: 10086,
				},
			},
		}
		err = initRegistry.Register(context.TODO(), service1)
		if err != nil {
			t.Fatal(err)
		}
	}()

	// 增加一个新服务
	go func() {
		time.Sleep(time.Second * 10)
		serviceName3 := "comment_service_1"
		service3 := &registry.Service{
			Name: serviceName3,
			Nodes: []*registry.Node{
				{
					Id:   0,
					Ip:   "127.0.0.1",
					Port: 68001,
				},
				{
					Id:   1,
					Ip:   "127.0.0.2",
					Port: 68001,
				},
			},
		}
		err = initRegistry.Register(context.TODO(), service3)
		if err != nil {
			t.Fatal(err)
		}
	}()

	// 等待服务全部注册完成
	time.Sleep(time.Second * 5)
	ticker := time.After(time.Second * 100)
	for {
		select {
		case <-ticker:
			goto end
		default:
			getService, err := initRegistry.GetService(context.TODO(), serviceName0)
			if err != nil {
				t.Fatal(err)
			}
			for i, node := range getService.Nodes {
				fmt.Printf("service0:%s node:%d,%+v ", getService.Name, i, node)
			}
			fmt.Println()
			time.Sleep(time.Second)
		}
	}
end:
}

func BenchmarkEtcdRegistry_GetService(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = initRegistry.GetService(context.TODO(), serviceName0)
	}
}
