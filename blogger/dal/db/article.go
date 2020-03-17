package db

import (
	"fmt"
	"github.com/summerKK/go-code-snippet-library/blogger/model"
)

func ArticleInsert(article *model.ArticleDetail) (articleId int64, err error) {
	if article == nil {
		return 0, fmt.Errorf("invalid params")
	}
	sqlStr := `insert into article 
    (category_id,content,title, username,summary)
    values (?,?,?,?,?)
    `
	result, err := Db.Exec(sqlStr, article.CategoryId, article.Content, article.Title, article.Username, article.Summary)
	if err != nil {
		return 0, err
	}

	articleId, _ = result.LastInsertId()
	return
}

func ArticleList(page int, pageSize int) (list []*model.ArticleInfo, err error) {
	if page < 0 || pageSize < 0 {
		return nil, fmt.Errorf("invalid params")
	}
	sqlStr := `select * from article where status = 1 order by create_time desc limit ?,?`
	err = Db.Unsafe().Select(&list, sqlStr, page, pageSize)
	if err != nil {
		return nil, err
	}
	return
}

func ArticleInfo(articleId int64) (articleInfo *model.ArticleDetail, err error) {
	articleInfo = &model.ArticleDetail{}
	sqlStr := "select * from article where  id = ?"
	err = Db.Unsafe().Get(articleInfo, sqlStr, articleId)
	return
}

func RelatedArticleList(articleId int64) (list []*model.RelatedArticle, err error) {
	sql := "select category_id from article where id = ?"
	var categortId int64
	err = Db.Get(&categortId, sql, articleId)
	if err != nil {
		return
	}
	sql = "select id,title from article where category_id = ? and id != ?"
	err = Db.Select(&list, sql, categortId, articleId)
	return
}

func PrevArticle(articleId int64) (article *model.RelatedArticle, err error) {
	article = &model.RelatedArticle{
		ArticleId: -1,
	}
	sql := "select id,title from article where id < ? order by id desc limit 1"
	err = Db.Get(article, sql, articleId)
	return
}

func NextArticle(articleId int64) (article *model.RelatedArticle, err error) {
	article = &model.RelatedArticle{
		ArticleId: -1,
	}
	sql := "select id,title from article where id > ? order by id asc limit 1"
	err = Db.Get(article, sql, articleId)
	return
}
