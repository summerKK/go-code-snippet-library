package upload

import (
	"io"
	"mime/multipart"
	"os"
	"path"
	"strings"

	"github.com/summerKK/go-code-snippet-library/blog-service/global"
	"github.com/summerKK/go-code-snippet-library/blog-service/pkg/util"
)

type FileType int

const (
	FileImage FileType = iota + 1
)

func GetFileName(name string) string {
	ext := GetFileExt(name)
	fileName := strings.TrimSuffix(name, ext)
	fileName = util.EncodeMd5(fileName)

	return fileName + ext
}

func GetFileExt(name string) string {
	return path.Ext(name)
}

func GetSavePath() string {
	return global.AppSetting.UploadSavePath
}

func CheckSavePath(dst string) bool {
	_, err := os.Stat(dst)

	return os.IsNotExist(err)
}

func CheckContainExt(t FileType, name string) bool {
	ext := GetFileExt(name)
	switch t {
	case FileImage:
		for _, configExt := range global.AppSetting.UploadImageAllowExts {
			if strings.ToLower(configExt) == strings.ToLower(ext) {
				return true
			}
		}
	}

	return false
}

// 检查文件是否超过最大字节
func CheckMaxSize(t FileType, name string) bool {
	fileInfo, err := os.Stat(name)
	if err != nil {
		return false
	}
	switch t {
	case FileImage:
		if fileInfo.Size() <= (global.AppSetting.UploadImageMaxSize << 20) {
			return true
		}
	}

	return false
}

func CheckPermission(dst string) bool {
	_, err := os.Stat(dst)

	return os.IsPermission(err)
}

func CreateSavePath(dst string, perm os.FileMode) error {
	return os.MkdirAll(dst, perm)
}

func SaveFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)

	return err
}
