package common

import "time"

const (
	QuestionStatusPending = iota + 1
)

type Question struct {
	Id         int64     `json:"id" db:"id"`
	QuestionId int64     `json:"question_id" db:"question_id"`
	Caption    string    `json:"caption" db:"caption"`
	Content    string    `json:"content" db:"content"`
	AuthorId   int64     `json:"author_id" db:"author_id"`
	CategoryId int64     `json:"category_id" db:"category_id"`
	CreateTime time.Time `json:"create_time" db:"create_time"`
	UpdateTime time.Time `json:"update_time" db:"update_time"`
	Status     int       `json:"status" db:"status"`
}
