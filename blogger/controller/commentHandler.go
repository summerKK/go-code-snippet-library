package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/summerKK/go-code-snippet-library/blogger/logic"
	"net/http"
	"strconv"
)

func CommentSubmit(c *gin.Context) {
	content := c.PostForm("comment")
	username := c.PostForm("username")
	articleIdStr := c.Query("article_id")

	articleId, err := strconv.ParseInt(articleIdStr, 10, 64)
	if err != nil {
		fmt.Printf("(0)comment insert got error:%v\n", err)
		c.HTML(http.StatusInternalServerError, "views/500.html", nil)
		return
	}
	_, err = logic.CommentInsert(content, username, articleId)
	if err != nil {
		fmt.Printf("(1)article info got error:%v\n", err)
		c.HTML(http.StatusInternalServerError, "views/500.html", nil)
		return
	}

	c.Redirect(http.StatusFound, "/article/detail/?article_id="+articleIdStr)
}
