package common

import "time"

type UserInfo struct {
	Id         int64     `json:"id" db:"id"`
	UserId     uint64    `json:"user_id" db:"user_id"`
	User       string    `json:"user" db:"username"`
	Nickname   string    `json:"nickname" db:"nickname"`
	Sex        int       `json:"sex" db:"sex"`
	Email      string    `json:"email" db:"email"`
	Password   string    `json:"password" db:"password"`
	CreateTime time.Time `json:"create_time" db:"create_time"`
	UpdateTime time.Time `json:"update_time" db:"update_time"`
}

const (
	UserSexMan   = 1
	UserSexWoman = 2
)
