package db

import "github.com/summerKK/go-code-snippet-library/blogger/model"

func CommentInsert(comment *model.Comment) (commentId int64, err error) {
	sql := "insert into comment(content,username,article_id) values (?,?,?)"
	result, err := Db.Exec(sql, comment.Content, comment.Username, comment.ArticleId)
	if err != nil {
		return
	}
	commentId, _ = result.LastInsertId()
	return
}

func CommentListByArticleId(articleId int64) (list []*model.Comment, err error) {
	sql := "select * from comment where article_id = ?"
	err = Db.Select(&list, sql, articleId)
	return
}
