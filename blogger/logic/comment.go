package logic

import (
	"summer/blogger/dal/db"
	"summer/blogger/model"
)

func CommentInsert(content, username string, articleId int64) (commentId int64, err error) {
	comment := &model.Comment{
		Content:   content,
		Username:  username,
		ArticleId: articleId,
	}
	commentId, err = db.CommentInsert(comment)
	return
}

func GetArticleCommentList(articleId int64) (list []*model.Comment, err error) {
	list, err = db.CommentListByArticleId(articleId)
	return
}
