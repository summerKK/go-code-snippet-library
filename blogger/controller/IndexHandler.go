package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"summer/blogger/logic"
)

func IndexHandler(c *gin.Context) {
	list, err := logic.GetArticleRecordList(0, 15)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "views/500.html", nil)
		return
	}

	c.HTML(http.StatusOK, "views/index.html", list)
}
