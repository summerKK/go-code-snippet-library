package account

import (
	"github.com/summerKK/go-code-snippet-library/spark/common"
	"github.com/summerKK/go-code-snippet-library/spark/dal/db/account"
	idgen "github.com/summerKK/go-code-snippet-library/spark/id-gen"
	"github.com/summerKK/go-code-snippet-library/spark/util"
)

const (
	UserSlat = "hpcNyjqIeAcWCGzYPQVQttqQKev9w4Rd"
)

func Register(user *common.UserInfo) (err error) {
	// 生成用户Id
	id, err := idgen.GenId()
	// 密码加盐
	if err != nil {
		return
	}
	user.UserId = id
	user.Password = util.Slat([]byte(user.Password + UserSlat))
	err = account.RegisterUser(user)
	return
}

func Login(user *common.UserInfo) (userInfo *common.UserInfo, err error) {
	user.Password = util.Slat([]byte(user.Password + UserSlat))
	userInfo, err = account.Login(user)
	return
}
