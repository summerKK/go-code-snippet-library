package question

import (
	"github.com/summerKK/go-code-snippet-library/spark/common"
	"github.com/summerKK/go-code-snippet-library/spark/dal/db"
	"github.com/summerKK/go-code-snippet-library/spark/logger"
)

func Save(question *common.Question) (err error) {
	sql := "insert into question( question_id, caption, content, author_id, category_id, status) VALUES (?,?,?,?,?,?)"
	_, err = db.Db.Query(sql, question.QuestionId, question.Caption, question.Content, question.AuthorId, question.CategoryId, question.Status)
	if err != nil {
		logger.Logger.Debug("insert into question get error:%v", err)
	}
	return
}

func List(categoryId, perPage, pageSize int64) (list []*common.QuestionDetail, err error) {
	sql := "select question.*,user.username as author_name from question left join user on question.author_id = user.user_id  where category_id = ? and status = 1 limit ?,?"
	err = db.Db.Select(&list, sql, categoryId, (pageSize-1)*perPage, perPage)
	if err != nil {
		logger.Logger.Debug("select qestion row get error:%v", err)
	}
	return
}
