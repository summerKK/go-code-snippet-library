package category

import (
	"github.com/summerKK/go-code-snippet-library/spark/common"
	"github.com/summerKK/go-code-snippet-library/spark/dal/db/category"
)

func List() (list []*common.Category, err error) {
	return category.List()
}
