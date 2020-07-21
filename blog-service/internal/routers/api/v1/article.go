package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/summerKK/go-code-snippet-library/blog-service/global"
	"github.com/summerKK/go-code-snippet-library/blog-service/internal/service"
	"github.com/summerKK/go-code-snippet-library/blog-service/pkg/app"
	"github.com/summerKK/go-code-snippet-library/blog-service/pkg/errcode"
)

type Article struct {
}

func NewArticle() Article {
	return Article{}
}

// @Summary 获取单个文章
// @Produce json
// @Param id path int true "文章ID"
// @Success 200 {object} model.Article "成功"
// @Failure 400 {object} errcode.Error "请求错误"
// @Failure 500 {object} errcode.Error "内部错误"
// @Router /api/v1/articles/{id} [get]
func (a Article) Get(c *gin.Context) {
	request := service.ArticleRequest{}
	response := app.NewResponse(c)
	ok, errors := app.BindAndValid(c, &request)
	if ok {
		global.Logger.Errorf("app.bindAndValid error:%v", errors)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errors.Errors()...))
		return
	}

	svc := service.New(c.Request.Context())
	article, err := svc.GetArticle(&request)
	if err != nil {
		global.Logger.Errorf("svc.GetArticle error:%v", err)
		response.ToErrorResponse(errcode.ErrorGetArticleFail)
		return
	}

	response.ToResponse(article)

	return
}

// @Summary 获取多个文章
// @Produce json
// @Param title query string false "文章名称"
// @Param tag_id query int false "标签ID"
// @Param state query int false "状态"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} model.ArticleSwagger "成功"
// @Failure 400 {object} errcode.Error "请求错误"
// @Failure 500 {object} errcode.Error "内部错误"
// @Router /api/v1/articles [get]
func (a Article) List(c *gin.Context) {
	request := service.ArticleListRequest{}
	response := app.NewResponse(c)
	ok, errors := app.BindAndValid(c, &request)
	if ok {
		global.Logger.Errorf("app.bindAndValid error:%v", errors)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errors.Errors()...))
		return
	}

	svc := service.New(c.Request.Context())
	pager := &app.Pager{
		PageSize: app.GetPageSize(c),
		Page:     app.GetPage(c),
	}
	list, totalRow, err := svc.GetArticleList(&request, pager)
	if err != nil {
		global.Logger.Errorf("svc.GetArticleList error:%v", err)
		response.ToErrorResponse(errcode.ErrorGetArticlesFail)
		return
	}
	pager.TotalRows = totalRow

	response.ToResponseList(list, totalRow)

	return
}

// @Summary 创建文章
// @Produce json
// @Param tag_id body string true "标签ID"
// @Param title body string true "文章标题"
// @Param desc body string false "文章简述"
// @Param cover_image_url body string true "封面图片地址"
// @Param content body string true "文章内容"
// @Param created_by body int true "创建者"
// @Param state body int false "状态"
// @Success 200 {object} model.Article "成功"
// @Failure 400 {object} errcode.Error "请求错误"
// @Failure 500 {object} errcode.Error "内部错误"
// @Router /api/v1/articles [post]
func (a Article) Create(c *gin.Context) {
	request := &service.CreateArticleRequest{}
	response := app.NewResponse(c)
	ok, errors := app.BindAndValid(c, request)
	if ok {
		global.Logger.Errorf("app.BindAndValid error:%v", errors)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errors.Errors()...))
		return
	}

	svc := service.New(c.Request.Context())
	err := svc.CreateArticle(request)
	if err != nil {
		global.Logger.Errorf("svc.CreateArticle error:%v", err)
		response.ToErrorResponse(errcode.ErrorCreateArticleFail)
		return
	}

	response.ToResponse(gin.H{})

	return
}

// @Summary 更新文章
// @Produce json
// @Param tag_id body string false "标签ID"
// @Param title body string false "文章标题"
// @Param desc body string false "文章简述"
// @Param cover_image_url body string false "封面图片地址"
// @Param content body string false "文章内容"
// @Param modified_by body string true "修改者"
// @Success 200 {object} model.Article "成功"
// @Failure 400 {object} errcode.Error "请求错误"
// @Failure 500 {object} errcode.Error "内部错误"
// @Router /api/v1/articles/{id} [put]
func (a Article) Update(c *gin.Context) {
	request := &service.UpdateArticleRequest{}
	response := app.NewResponse(c)
	ok, errors := app.BindAndValid(c, request)
	if ok {
		global.Logger.Errorf("app.BindAndValid error:%v", errors)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errors.Errors()...))
		return
	}

	svc := service.New(c.Request.Context())
	err := svc.UpdateArticle(request)
	if err != nil {
		global.Logger.Errorf("svc.UpdateArticle error:%v", err)
		response.ToErrorResponse(errcode.ErrorUpdateArticleFail)
		return
	}

	response.ToResponse(gin.H{})

	return
}

// @Summary 删除文章
// @Produce  json
// @Param id path int true "文章ID"
// @Success 200 {string} string "成功"
// @Failure 400 {object} errcode.Error "请求错误"
// @Failure 500 {object} errcode.Error "内部错误"
// @Router /api/v1/articles/{id} [delete]
func (a Article) Delete(c *gin.Context) {
	request := &service.DeleteArticleRequest{}
	response := app.NewResponse(c)
	ok, errors := app.BindAndValid(c, request)
	if ok {
		global.Logger.Errorf("app.BindAndValid error:%v", errors)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errors.Errors()...))
		return
	}

	svc := service.New(c.Request.Context())
	err := svc.DeleteArticle(request)
	if err != nil {
		global.Logger.Errorf("svc.DeleteArticle error:%v", err)
		response.ToErrorResponse(errcode.ErrorDeleteArticleFail)
		return
	}

	response.ToResponse(gin.H{})

	return
}
