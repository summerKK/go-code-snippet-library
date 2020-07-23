package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/summerKK/go-code-snippet-library/blog-service/global"
	"github.com/summerKK/go-code-snippet-library/blog-service/internal/service"
	"github.com/summerKK/go-code-snippet-library/blog-service/pkg/app"
	"github.com/summerKK/go-code-snippet-library/blog-service/pkg/convert"
	"github.com/summerKK/go-code-snippet-library/blog-service/pkg/errcode"
)

type Tag struct {
}

func NewTag() Tag {
	return Tag{}
}

// @Summary 获取多个标签
// @Produce json
// @Param name query string false "标签名称" maxlength(100)
// @Param state query int false "标签名称" Enums(0,1) default(1)
// @Param page query int false "页码"
// @Param page_size query int false "每页个数"
// @Success 200 {object} model.TagSwagger "成功"
// @Failure 400 {object} errcode.Error "请求错误"
// @Failure 500 {object} errcode.Error "内部错误"
// @Router /api/v1/tags [get]
func (t Tag) List(c *gin.Context) {
	params := service.TagListRequest{}
	response := app.NewResponse(c)
	ok, errors := app.BindAndValid(c, &params)
	if ok {
		global.Logger.Errorf(c, "app.BindAndValid error:%v", errors)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errors.Errors()...))
		return
	}

	svc := service.New(c.Request.Context())
	countTag, err := svc.CountTag(&service.CountTagRequest{Name: params.Name, State: params.State})
	if err != nil {
		global.Logger.Errorf(c, "svc.CountTag error:%v", err)
		response.ToErrorResponse(errcode.ErrorCountTagFail)
		return
	}

	pager := app.Pager{
		Page:      app.GetPage(c),
		PageSize:  app.GetPageSize(c),
		TotalRows: countTag,
	}

	list, err := svc.GetTagList(&params, &pager)
	if err != nil {
		global.Logger.Errorf(c, "svc.GetTagList error:%v", err)
		response.ToErrorResponse(errcode.ErrorGetTagListFail)
		return
	}

	response.ToResponseList(list, countTag)

	return
}

// @Summary 新增标签
// @Produce  json
// @Param name body string true "标签名称" minlength(3) maxlength(100)
// @Param state body int false "状态" Enums(0, 1) default(1)
// @Param created_by body string false "创建者" minlength(3) maxlength(100)
// @Success 200 {object} model.Tag "成功"
// @Failure 400 {object} errcode.Error "请求错误"
// @Failure 500 {object} errcode.Error "内部错误"
// @Router /api/v1/tags [post]
func (t Tag) Create(c *gin.Context) {
	params := service.CreateTagRequest{}
	response := app.NewResponse(c)
	ok, errors := app.BindAndValid(c, &params)
	if ok {
		global.Logger.Errorf(c, "app.BindAndValid error:%v", errors)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errors.Errors()...))
		return
	}

	svc := service.New(c.Request.Context())
	err := svc.CreateTag(&params)
	if err != nil {
		global.Logger.Errorf(c, "svc.CreateTag error:%v", err)
		response.ToErrorResponse(errcode.ErrorCreateTagFail)
		return
	}

	response.ToResponse(gin.H{})

	return
}

// @Summary 更新标签
// @Produce  json
// @Param id path int true "标签ID"
// @Param name body string false "标签名称" minlength(3) maxlength(100)
// @Param state body int false "状态" Enums(0, 1) default(1)
// @Param modified_by body string true "修改者" minlength(3) maxlength(100)
// @Success 200 {array} model.Tag "成功"
// @Failure 400 {object} errcode.Error "请求错误"
// @Failure 500 {object} errcode.Error "内部错误"
// @Router /api/v1/tags/{id} [put]
func (t Tag) Update(c *gin.Context) {
	params := service.UpdateTagRequest{
		ID: convert.StrTo(c.Param("id")).MustUInt32(),
	}
	response := app.NewResponse(c)
	ok, errors := app.BindAndValid(c, &params)
	if ok {
		global.Logger.Errorf(c, "app.BidAndValid error:%v", errors)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errors.Errors()...))
		return
	}

	svc := service.New(c.Request.Context())
	err := svc.UpdateTag(&params)
	if err != nil {
		global.Logger.Errorf(c, "svc.UpdateTag error:%v", err)
		response.ToErrorResponse(errcode.ErrorUpdateTagFail)
		return
	}

	response.ToResponse(gin.H{})

	return
}

// @Summary 删除标签
// @Produce  json
// @Param id path int true "标签ID"
// @Success 200 {string} string "成功"
// @Failure 400 {object} errcode.Error "请求错误"
// @Failure 500 {object} errcode.Error "内部错误"
// @Router /api/v1/tags/{id} [delete]
func (t Tag) Delete(c *gin.Context) {
	params := service.DeleteTagRequest{}
	response := app.NewResponse(c)
	ok, errors := app.BindAndValid(c, &params)
	if ok {
		global.Logger.Errorf(c, "app.BindAndValid error:%v", errors)
		response.ToErrorResponse(errcode.ErrorDeleteTagFail.WithDetails(errors.Errors()...))
		return
	}

	svc := service.New(c.Request.Context())
	err := svc.DeleteTag(&params)
	if err != nil {
		global.Logger.Errorf(c, "svc.DeleteTag error:%v", err)
		response.ToErrorResponse(errcode.ErrorDeleteTagFail)
		return
	}

	response.ToResponse(gin.H{})

	return
}
