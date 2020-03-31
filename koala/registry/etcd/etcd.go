package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/summerKK/go-code-snippet-library/koala/logger"
	"github.com/summerKK/go-code-snippet-library/koala/registry"
	"go.etcd.io/etcd/clientv3"
	"path"
	"time"
)

const (
	maxServiceLen = 8
)

type etcdRegistry struct {
	options            *registry.Options
	client             *clientv3.Client
	serviceCh          chan *registry.Service
	registerServiceMap map[string]*registerService
}

type registerService struct {
	id          clientv3.LeaseID
	service     *registry.Service
	registered  bool
	keepAliveCh <-chan *clientv3.LeaseKeepAliveResponse
}

var (
	etcdRegistryEntry = etcdRegistry{
		serviceCh:          make(chan *registry.Service, maxServiceLen),
		registerServiceMap: make(map[string]*registerService, 8),
	}
)

// 自动注册到注册中心组建
func init() {
	registry.RegisterRegistry(&etcdRegistryEntry)
	go etcdRegistryEntry.run()
}

func (e *etcdRegistry) GetName() (name string) {
	return "etcd"
}

func (e *etcdRegistry) Register(ctx context.Context, service *registry.Service) (err error) {
	select {
	case e.serviceCh <- service:
		return
	default:
		err = fmt.Errorf("register service faild")
		return
	}
}

func (e *etcdRegistry) UnRegister(ctx context.Context, service *registry.Service) (err error) {
	panic("implement me")
}

func (e *etcdRegistry) Init(ctx context.Context, options ...registry.Option) (err error) {
	e.options = &registry.Options{}
	for _, option := range options {
		option(e.options)
	}
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   e.options.Addrs,
		DialTimeout: e.options.Timeout,
	})
	if err != nil {
		return
	}

	e.client = cli
	return
}

func (e *etcdRegistry) run() {
	for {
		select {
		case service := <-e.serviceCh:
			// 查看服务是否已经注册过
			_, ok := e.registerServiceMap[service.Name]
			if !ok {
				rs := &registerService{
					service: service,
				}
				// 添加到map,然后在`registerOrKeepAlive`统一注册
				e.registerServiceMap[service.Name] = rs
			}
		default:
			e.registerOrKeepAlive()
			time.Sleep(time.Millisecond * 500)
		}
	}

}

func (e *etcdRegistry) registerOrKeepAlive() {
	for _, rs := range e.registerServiceMap {
		// 已经注册过,检查服务是否有异常
		if rs.registered {
			err := e.serviceKeepAlive(rs)
			if err != nil {
				logger.Logger.Warn(err)
			}
			continue
		}

		err := e.registerService(rs)
		if err != nil {
			logger.Logger.Warn(err)
		}
	}
}

func (e *etcdRegistry) registerService(rs *registerService) (err error) {
	// 续期
	resp, err := e.client.Grant(context.TODO(), e.options.HeartBet)
	if err != nil {
		return
	}
	rs.id = resp.ID

	for _, node := range rs.service.Nodes {
		tmp := &registry.Service{
			Name: rs.service.Name,
			Nodes: []*registry.Node{
				node,
			},
		}
		key := e.serviceNodePath(tmp)
		data, err := json.Marshal(tmp)
		// 出错留到下次注册
		if err != nil {
			continue
		}
		_, err = e.client.Put(context.TODO(), key, string(data), clientv3.WithLease(resp.ID))
		if err != nil {
			continue
		}
	}

	// key永久保持
	alive, err := e.client.KeepAlive(context.TODO(), resp.ID)
	if err != nil {
		return
	}

	rs.keepAliveCh = alive
	rs.registered = true

	return
}

func (e *etcdRegistry) serviceKeepAlive(rs *registerService) (err error) {
	select {
	case resp := <-rs.keepAliveCh:
		if resp == nil {
			rs.registered = false
			return
		}
		logger.Logger.Infof("service:%s, ttl:%v", rs.service.Name, resp.TTL)
	}
	return
}

func (e *etcdRegistry) serviceNodePath(service *registry.Service) string {
	nodeIp := fmt.Sprintf("%s:%d", service.Nodes[0].Ip, service.Nodes[0].Port)
	return path.Join(e.options.RegistryPath, service.Name, nodeIp)
}
