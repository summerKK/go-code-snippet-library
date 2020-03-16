package session

import (
	"fmt"
)

var Session ISessionMgr

func Init(provider string, addr string, options ...string) (err error) {
	switch provider {
	case "redis":
		Session, err = NewRedisSessionMgr(addr, options...)
		if err != nil {
			err = fmt.Errorf("init redis session failed,error:%v\n", err)
		}
	case "memory":
		Session, err = NewMemSessionMgr(addr, options...)
		if err != nil {
			err = fmt.Errorf("init memory session failed,error:%v\n", err)
		}
	default:
		err = fmt.Errorf("unsupported provider:%v\n", err)
	}
	return
}
