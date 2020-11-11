package gin

import (
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"sort"
)

const (
	AUTH_USER_KEY = "user"
)

type BasicAuthPair struct {
	Code string
	User string
}

type Accounts map[string]string

type Pairs []BasicAuthPair

func (p Pairs) Len() int {
	return len(p)
}

func (p Pairs) Less(i, j int) bool {
	// 使用升序排列
	return p[i].Code < p[j].Code
}

func (p Pairs) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func ProcessCredentials(accounts Accounts) (Pairs, error) {
	if len(accounts) == 0 {
		return nil, errors.New("Empty list of authorized credentials.")
	}

	pairs := make(Pairs, 0, len(accounts))
	for account, password := range accounts {
		if len(account) == 0 || len(password) == 0 {
			return nil, errors.New("User or Password empty.")
		}
		base := account + ":" + password
		code := "Basic " + base64.StdEncoding.EncodeToString([]byte(base))
		pairs = append(pairs, BasicAuthPair{Code: code, User: account})
	}
	// 排序后通过二分查找进行查找.时间复杂度为O(logN)
	sort.Sort(pairs)

	return pairs, nil
}

func SearchCredential(pairs Pairs, auth string) string {
	if len(auth) == 0 {
		return ""
	}

	index := sort.Search(len(pairs), func(i int) bool {
		// 升序的时候用 `>=`,降序的时候用 `<=`
		return pairs[i].Code >= auth
	})

	// subtle.ConstantTimeCompare([]byte(pairs[index].Code), []byte(auth)) == 1
	// see https://zhuanlan.zhihu.com/p/143270224
	if index < len(pairs) && subtle.ConstantTimeCompare([]byte(pairs[index].Code), []byte(auth)) == 1 {
		return pairs[index].User
	} else {
		return ""
	}
}

func BasicAuth(accounts Accounts) HandlerFunc {
	pairs, err := ProcessCredentials(accounts)
	if err != nil {
		panic(err)
	}
	return func(c *Context) {
		user := SearchCredential(pairs, c.Req.Header.Get("Authorization"))
		if len(user) == 0 {
			c.Writer.Header().Set("WWW-Authenticate", "Basic realm=\"Authorization Required\"")
			c.Fail(401, errors.New("Unauthorized"))
		} else {
			c.Set(AUTH_USER_KEY, user)
		}

		c.Next()
	}
}
