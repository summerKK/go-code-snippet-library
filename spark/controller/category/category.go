package category

import (
	"github.com/gin-gonic/gin"
	"github.com/summerKK/go-code-snippet-library/spark/service/category"
	"github.com/summerKK/go-code-snippet-library/spark/util"
)

func List(c *gin.Context) {
	list, err := category.List()
	if err != nil {
		util.ResponseErr(c, &util.CodeInfo{Code: util.ErrServiceBusy})
		return
	}
	util.ResponseSuc(c, list, nil)
}
