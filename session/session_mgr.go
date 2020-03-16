package session

type ISessionMgr interface {
	Init(addr string, options ...string) (err error)
	Get(sessionId string) (session ISession, err error)
	Create() (session ISession, err error)
}
