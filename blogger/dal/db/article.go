package db

import (
	"fmt"
	"summer/blogger/model"
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
