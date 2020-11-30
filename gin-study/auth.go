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

type Accounts map[string]string

type authPair struct {
	Value string
	User  string
}

type AuthPairs []authPair

func (p AuthPairs) Len() int {
	return len(p)
}

func (p AuthPairs) Less(i, j int) bool {
	// 使用升序排列
	return p[i].Value < p[j].Value
}

func (p AuthPairs) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func ProcessAccounts(accounts Accounts) (AuthPairs, error) {
	if len(accounts) == 0 {
		return nil, errors.New("Empty list of authorized credentials.")
	}

	pairs := make(AuthPairs, 0, len(accounts))
	for account, password := range accounts {
		if len(account) == 0 || len(password) == 0 {
			return nil, errors.New("User or Password empty.")
		}
		base := account + ":" + password
		code := "Basic " + base64.StdEncoding.EncodeToString([]byte(base))
		pairs = append(pairs, authPair{Value: code, User: account})
	}
	// 排序后通过二分查找进行查找.时间复杂度为O(logN)
	sort.Sort(pairs)

	return pairs, nil
}

func SearchAccount(pairs AuthPairs, auth string) (string, bool) {
	if len(auth) == 0 {
		return "", false
	}

	index := sort.Search(len(pairs), func(i int) bool {
		// 升序的时候用 `>=`,降序的时候用 `<=`
		return pairs[i].Value >= auth
	})

	if index < len(pairs) && secureCompare(pairs[index].Value, auth) {
		return pairs[index].User, true
	} else {
		return "", false
	}
}

// subtle.ConstantTimeCompare([]byte(pairs[index].Code), []byte(auth)) == 1
// see https://zhuanlan.zhihu.com/p/143270224
func secureCompare(given, actual string) bool {
	if subtle.ConstantTimeEq(int32(len(given)), int32(len(actual))) == 1 {
		return subtle.ConstantTimeCompare([]byte(given), []byte(actual)) == 1
	} else {
		/* Securely compare actual to itself to keep constant time, but always return false */
		return false
	}
}

func BasicAuth(accounts Accounts) HandlerFunc {
	pairs, err := ProcessAccounts(accounts)
	if err != nil {
		panic(err)
	}
	return func(c *Context) {
		user, ok := SearchAccount(pairs, c.Request.Header.Get("Authorization"))
		if !ok {
			c.Writer.Header().Set("WWW-Authenticate", "Basic realm=\"Authorization Required\"")
			c.Fail(401, errors.New("Unauthorized"))
		} else {
			c.Set(AUTH_USER_KEY, user)
		}

		c.Next()
	}
}
