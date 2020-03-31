package registry

import "time"

// 参数
type Options struct {
	Addrs   []string
	Timeout time.Duration
	// 心跳检测时间
	HeartBet int64
	// /order:192.168.1.1:10086
	// /order:192.168.1.1:10087
	RegistryPath string
}

type Option func(o *Options)

func WithAddrs(addrs []string) Option {
	return func(o *Options) {
		o.Addrs = addrs
	}
}

func WithTimeout(time time.Duration) Option {
	return func(o *Options) {
		o.Timeout = time
	}
}
