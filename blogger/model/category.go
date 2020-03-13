package model

import "time"

type Category struct {
	CategoryId   int64     `db:id`
	CategoryName string    `db:category_name`
	CategoryNo   int       `db:category_no`
	CreateTime   time.Time `db:create_time`
}
