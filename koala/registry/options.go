package registry

import "time"

// 参数
type Options struct {
	Addrs   []string
	Timeout time.Duration
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
