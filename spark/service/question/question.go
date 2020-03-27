package question

import (
	"github.com/gin-gonic/gin"
	"github.com/summerKK/go-code-snippet-library/spark/common"
	"github.com/summerKK/go-code-snippet-library/spark/dal/db/question"
	idgen "github.com/summerKK/go-code-snippet-library/spark/id-gen"
	"github.com/summerKK/go-code-snippet-library/spark/logger"
	"github.com/summerKK/go-code-snippet-library/spark/middleware/account"
)

func Save(c *gin.Context, q *common.Question) (err error) {
	userId, err := account.GetUserId(c)
	if err != nil {
		logger.Logger.Debug("get user id failed:%v", err)
		return
	}
	q.AuthorId = userId
	id, err := idgen.GenId()
	if err != nil {
		logger.Logger.Debug("gen q id failed:%v", err)
		return
	}
	q.QuestionId = int64(id)
	q.Status = common.QuestionStatusPending
	err = question.Save(q)
	return
}
