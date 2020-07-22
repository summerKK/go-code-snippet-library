package api

import (
	"github.com/gin-gonic/gin"
	"github.com/summerKK/go-code-snippet-library/blog-service/internal/service"
	"github.com/summerKK/go-code-snippet-library/blog-service/pkg/app"
	"github.com/summerKK/go-code-snippet-library/blog-service/pkg/convert"
	"github.com/summerKK/go-code-snippet-library/blog-service/pkg/errcode"
	"github.com/summerKK/go-code-snippet-library/blog-service/pkg/upload"
)

// @Summary 上传文件
// @Produce json
// @Param file body string true "上传文件"
// @Param type body int true "上传文件类型 1图片"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} service.FileInfo "成功"
// @Failure 400 {object} errcode.Error "请求错误"
// @Failure 500 {object} errcode.Error "内部错误"
// @Router /api/upload/file [post]
func UploadFile(c *gin.Context) {
	response := app.NewResponse(c)
	file, header, err := c.Request.FormFile("file")
	fileType := convert.StrTo(c.PostForm("type")).MustInt()
	if err != nil {
		response.ToErrorResponse(errcode.ErrorFileUploadFail)
		return
	}
	if header == nil || fileType <= 0 {
		response.ToErrorResponse(errcode.InvalidParams)
		return
	}

	svc := service.New(c.Request.Context())
	fileInfo, err := svc.Upload(upload.FileType(fileType), file, header)
	if err != nil {
		response.ToErrorResponse(errcode.ErrorFileUploadFail.WithDetails(err.Error()))
		return
	}

	response.ToResponse(gin.H{
		"file_access_url": fileInfo.AccessUrl,
	})

	return
}
