package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/summerKK/go-code-snippet-library/koala/logger"
	"github.com/summerKK/go-code-snippet-library/koala/registry"
	"go.etcd.io/etcd/clientv3"
	"path"
	"sync"
	"sync/atomic"
	"time"
)

const (
	maxServiceNum     = 8
	syncCacheIntervel = 5
)

type etcdRegistry struct {
	options            *registry.Options
	client             *clientv3.Client
	serviceCh          chan *registry.Service
	registerServiceMap map[string]*registerService
	//  原子操作,添加value
	value atomic.Value
	sync.Mutex
}

type registerService struct {
	id          clientv3.LeaseID
	service     *registry.Service
	registered  bool
	keepAliveCh <-chan *clientv3.LeaseKeepAliveResponse
}

type allService struct {
	serviceMap map[string]*registry.Service
}

var (
	etcdRegistryEntry = etcdRegistry{
		serviceCh:          make(chan *registry.Service, maxServiceNum),
		registerServiceMap: make(map[string]*registerService, maxServiceNum),
	}
)

// 自动注册到注册中心组建
func init() {
	// 服务缓存
	s := &allService{serviceMap: make(map[string]*registry.Service, maxServiceNum)}
	etcdRegistryEntry.value.Store(s)

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
	ticker := time.Tick(time.Second * syncCacheIntervel)
	for {
		select {
		case service := <-e.serviceCh:
			// 查看服务是否已经注册过
			rds, ok := e.registerServiceMap[service.Name]
			if !ok {
				rs := &registerService{
					service: service,
				}
				// 添加到map,然后在`registerOrKeepAlive`统一注册
				e.registerServiceMap[service.Name] = rs
			} else {
				// service.Name已经存在,注册进来该服务的一个节点(需要考虑重复注册)
				reset := false
				for _, node0 := range service.Nodes {
					// 查看是否已经注册过了
					for _, node1 := range rds.service.Nodes {
						if fmt.Sprintf("%s_%d", node0.Ip, node0.Port) == fmt.Sprintf("%s_%d", node1.Ip, node1.Port) {
							goto dumplicate
						}
					}
					reset = true
					rds.service.Nodes = append(rds.service.Nodes, &registry.Node{
						Id:   len(rds.service.Nodes) + 1,
						Ip:   node0.Ip,
						Port: node0.Port,
					})
					logger.Logger.Infof("service:%s add new node %s:%d", rds.service.Name, node0.Ip, node0.Port)
				dumplicate:
				}
				if reset {
					// 服务增加新节点,重新注册
					rds.registered = false
				}
			}
		case <-ticker:
			e.syncServiceCache()
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
		logger.Logger.Infof("registry service key:%s", key)
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
		//logger.Logger.Infof("service:%s, ttl:%v", rs.service.Name, resp.TTL)
	}
	return
}

func (e *etcdRegistry) serviceNodePath(service *registry.Service) string {
	nodeIp := fmt.Sprintf("%s:%d", service.Nodes[0].Ip, service.Nodes[0].Port)
	return path.Join(e.options.RegistryPath, service.Name, nodeIp)
}

func (e *etcdRegistry) servicePath(name string) string {
	return path.Join(e.options.RegistryPath, name)
}

func (e *etcdRegistry) getServiceFromCache(name string) (service *registry.Service, exist bool) {
	// 这里是原子操作,不用加锁
	services := e.value.Load().(*allService)
	service, exist = services.serviceMap[name]
	return
}

func (e *etcdRegistry) GetService(ctx context.Context, name string) (service *registry.Service, err error) {
	name = e.servicePath(name)
	service, exist := e.getServiceFromCache(name)
	if exist {
		return
	}

	// 这里加锁,只允许一个线程从etcd加载配置.防止大量线程都从etcd里面拿数据造成雪崩
	e.Lock()
	defer e.Unlock()
	// 再次查看是否已经把配置信息加载到缓存里面去了
	service, exist = e.getServiceFromCache(name)
	if exist {
		return
	}

	response, err := e.client.Get(ctx, name, clientv3.WithPrefix())
	if err != nil {
		return
	}

	if response.Kvs == nil {
		return nil, fmt.Errorf("service: %s , empty node\n", name)
	}

	service = &registry.Service{
		Name: name,
	}

	for _, kv := range response.Kvs {
		var s registry.Service
		err = json.Unmarshal(kv.Value, &s)
		if err != nil {
			continue
		}

		for i, node := range s.Nodes {
			service.Nodes = append(service.Nodes, &registry.Node{
				Id:   i,
				Ip:   node.Ip,
				Port: node.Port,
			})
		}

		// 把服务保存在缓存中
		as := e.value.Load().(*allService)
		as.serviceMap[name] = service
		e.value.Store(as)
	}

	return
}

// 更新服务缓存信息,可能有新的服务注册.要把它们加到缓存里面去
func (e *etcdRegistry) syncServiceCache() {
	ctx := context.TODO()
	serviceCache := e.value.Load().(*allService)
	for name := range e.registerServiceMap {
		name = e.servicePath(name)
		response, err := e.client.Get(ctx, name, clientv3.WithPrefix())
		if err != nil {
			continue
		}

		var s registry.Service
		s.Name = name
		for _, kv := range response.Kvs {
			var tmp registry.Service
			err := json.Unmarshal(kv.Value, &tmp)
			if err != nil {
				continue
			}
			for i, node := range tmp.Nodes {
				s.Nodes = append(s.Nodes, &registry.Node{
					Id:   i,
					Ip:   node.Ip,
					Port: node.Port,
				})
			}
		}
		serviceCache.serviceMap[name] = &s
		e.value.Store(serviceCache)
	}
}
