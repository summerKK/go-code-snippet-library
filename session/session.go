package session

type ISession interface {
	Set(key string, value interface{}) (err error)
	Get(key string) (value interface{}, err error)
	Del(key string) (err error)
	Save() (err error)
	IsModify() bool
	Id() string
}

const (
	status_init = iota
	status_loading
	status_modify
)
