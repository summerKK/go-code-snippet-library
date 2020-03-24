package account

import (
	"github.com/summerKK/go-code-snippet-library/spark/common"
	"github.com/summerKK/go-code-snippet-library/spark/dal/db"
	"log"
)

func RegisterUser(user *common.UserInfo) (err error) {
	sql := "select count(*) from user where username = ?"
	var userCount int
	err = db.Db.Get(&userCount, sql, user.User)
	if err != nil {
		log.Printf("select user,got error:%v\n", err)
		return
	}
	if userCount > 0 {
		err = DbUserExists
		return
	}
	sql = "insert into user(user_id, username, nickname, password, email, sex) values (?,?,?,?,?,?)"
	_, err = db.Db.Exec(sql, user.UserId, user.User, user.Nickname, user.Password, user.Email, user.Sex)
	if err != nil {
		log.Printf("insert into user,got error:%v\n", err)
	}
	return
}
