package service

import (
	"errors"
	"fmt"
	"mime/multipart"
	"os"

	"github.com/summerKK/go-code-snippet-library/blog-service/global"
	"github.com/summerKK/go-code-snippet-library/blog-service/pkg/upload"
)

type FileInfo struct {
	Name      string `json:"name"`
	AccessUrl string `json:"access_url"`
}

// 文件上传
func (s Service) Upload(fileType upload.FileType, file multipart.File, header *multipart.FileHeader) (*FileInfo, error) {
	filename := upload.GetFileName(header.Filename)
	uploadSavePath := upload.GetSavePath()
	dst := uploadSavePath + "/" + filename
	if !upload.CheckContainExt(fileType, filename) {
		return nil, errors.New(fmt.Sprintf("不支持的文件类型:%s", upload.GetFileExt(filename)))
	}

	if upload.CheckSavePath(dst) {
		if err := upload.CreateSavePath(uploadSavePath, os.ModePerm); err != nil {
			return nil, errors.New("创建文件失败")
		}
	}

	if upload.CheckPermission(dst) {
		return nil, errors.New("文件权限不足")
	}

	if !upload.CheckMaxSize(fileType, header) {
		return nil, errors.New("文件超过最大上传尺寸")
	}

	if err := upload.SaveFile(header, dst); err != nil {
		return nil, errors.New("上传文件失败")
	}

	accessUrl := global.AppSetting.UploadServerUrl + "/" + filename

	return &FileInfo{
		Name:      filename,
		AccessUrl: accessUrl,
	}, nil

}
