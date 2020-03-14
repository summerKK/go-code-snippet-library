package model

type RelatedArticle struct {
	ArticleId int64  `db:"id"`
	Title     string `db:"title"`
}
