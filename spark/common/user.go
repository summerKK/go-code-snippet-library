package common

type UserInfo struct {
	Id       int64  `json:"id" db:"id"`
	UserId   uint64 `json:"user_id" db:"user_id"`
	User     string `json:"user" db:"username"`
	Nickname string `json:"nickname" db:"nickname"`
	Sex      int    `json:"sex" db:"sex"`
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
}
