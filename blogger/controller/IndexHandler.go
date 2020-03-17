package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/summerKK/go-code-snippet-library/blogger/logic"
	"net/http"
)

func IndexHandler(c *gin.Context) {
	list, err := logic.GetArticleRecordList(0, 15)
	if err != nil {
		fmt.Printf("index handler got error:%v\n", err)
		c.HTML(http.StatusInternalServerError, "views/500.html", nil)
		return
	}

	c.HTML(http.StatusOK, "views/index.html", list)
}
