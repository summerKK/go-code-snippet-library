package category

import (
	"github.com/summerKK/go-code-snippet-library/spark/common"
	"github.com/summerKK/go-code-snippet-library/spark/dal/db"
)

func List() (list []*common.Category, err error) {
	sql := "select * from spark.category"
	err = db.Db.Select(&list, sql)
	return
}
