package model

type User struct {
	Id       int64  `db:"id"`
	UserId   int64  `db:"user_id"`
	Username string `db:"username"`
	Nickname string `db:"nickname"`
	Password string `db:"password"`
}
