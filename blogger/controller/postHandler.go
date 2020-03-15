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

func ArticleInfo(c *gin.Context) {
	articleIdStr := c.Query("article_id")
	articleId, err := strconv.ParseInt(articleIdStr, 10, 64)
	if err != nil {
		fmt.Printf("(0)article info got error:%v\n", err)
		c.HTML(http.StatusInternalServerError, "views/500.html", nil)
		return
	}
	articleInfo, err := logic.ArticleInfo(articleId)
	if err != nil {
		fmt.Printf("(1)article info got error:%v\n", err)
		c.HTML(http.StatusInternalServerError, "views/500.html", nil)
		return
	}
	category, _ := logic.GetCategoryById(articleInfo.CategoryId)
	articleInfo.Category = *category

	relatedArticles, err := logic.RelatedArticleListById(articleId)
	if err != nil {
		fmt.Printf("(2) article info got error:%v\n", err)
	}

	prevArticle, err := logic.GetPrevArticleById(articleId)
	if err != nil {
		fmt.Printf("(3) article info got error:%v\n", err)
	}

	nextArticle, err := logic.GetNextArticleById(articleId)
	if err != nil {
		fmt.Printf("(3) article info got error:%v\n", err)
	}

	categories, err := logic.GetCategoryList()
	if err != nil {
		fmt.Printf("(4) article info got error:%v\n", err)
	}

	commentList, err := logic.GetArticleCommentList(articleId)
	if err != nil {
		fmt.Printf("(5) article info got error:%v\n", err)
	}

	var data map[string]interface{} = make(map[string]interface{}, 10)
	data["detail"] = articleInfo
	data["related_article"] = relatedArticles
	data["prev_article"] = prevArticle
	data["next_article"] = nextArticle
	data["categories"] = categories
	data["comment"] = commentList
	data["article_id"] = articleId

	c.HTML(http.StatusOK, "views/detail.html", data)
}
