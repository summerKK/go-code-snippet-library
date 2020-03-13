package model

import "time"

type ArticleInfo struct {
	Id           int64     `db:"id"`
	CategoryId   int64     `db:"category_id"`
	Title        string    `db:"title"`
	Summary      string    `db:"summary"`
	ViewCount    int64     `db:"view_count"`
	CommentCount int64     `db:"comment_count"`
	Username     string    `db:"username"`
	CreateTime   time.Time `db:"create_time"`
}

type ArticleDetail struct {
	ArticleInfo
	Content string `db:"content"`
	Category
}

type ArticleRecord struct {
	ArticleInfo
	Category
}
