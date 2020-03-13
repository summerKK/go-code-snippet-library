package model

import "time"

type Leave struct {
	Id         int64     `db:"id"`
	Username   string    `db:"username"`
	Email      string    `db:"email"`
	Content    string    `db:"content"`
	CreateTime time.Time `db:"create_time"`
}
