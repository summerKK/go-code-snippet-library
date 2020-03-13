package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"summer/blogger/logic"
)

func ArticleCreate(c *gin.Context) {
	list, err := logic.GetCategoryList()
	if err != nil {
		fmt.Printf("article create got error:%v\n", err)
		c.HTML(http.StatusInternalServerError, "views/500.html", nil)
		return
	}
	c.HTML(http.StatusOK, "views/post_article.html", list)
}

func ArticleSubmit(c *gin.Context) {
	username := c.PostForm("author")
	title := c.PostForm("title")
	categoryIdStr := c.PostForm("category_id")
	content := c.PostForm("content")
	categoryId, err := strconv.ParseInt(categoryIdStr, 10, 64)
	if err != nil {
		fmt.Printf("(0)article submit got error:%v\n", err)
		c.HTML(http.StatusInternalServerError, "views/500.html", nil)
		return
	}

	err = logic.ArticleInsert(username, title, content, categoryId)
	if err != nil {
		fmt.Printf("(1)article submit got error:%v\n", err)
		c.HTML(http.StatusInternalServerError, "views/500.html", nil)
		return
	}

	c.Redirect(http.StatusFound, "/")
}
