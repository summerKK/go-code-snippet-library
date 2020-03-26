package common

import "time"

type Category struct {
	Id           int64     `db:"id" json:"id"`
	CategoryId   int64     `db:"category_id" json:"category_id"`
	CategoryName string    `db:"category_name" json:"category_name"`
	CreateTime   time.Time `db:"create_time" json:"create_time"`
	UpdateTime   time.Time `db:"update_time" json:"update_time"`
}
