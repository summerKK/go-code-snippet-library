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
