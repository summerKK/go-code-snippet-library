package account

import "github.com/summerKK/go-code-snippet-library/session"

func InitSession() (err error) {
	return session.Init("redis", "127.0.0.1:6379")
}
