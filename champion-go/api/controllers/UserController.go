package controllers

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/summerKK/go-code-snippet-library/champion-go/api/auth"
	"github.com/summerKK/go-code-snippet-library/champion-go/api/models"
	"github.com/summerKK/go-code-snippet-library/champion-go/api/security"
	"github.com/summerKK/go-code-snippet-library/champion-go/api/utils/fileformat"
	"github.com/summerKK/go-code-snippet-library/champion-go/api/utils/formatError"
	"golang.org/x/crypto/bcrypt"
)

func (s *Server) CreateUser(c *gin.Context) {

	errList := make(map[string]string)
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		errList["Invalid_body"] = "Unable to get request"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	user := &models.User{}
	err = json.Unmarshal(body, user)
	if err != nil {
		errList["Unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	user.Prepare()
	errMsgMap := user.Validate(models.USER_TYPE_DEFAULT)
	if len(errMsgMap) > 0 {
		errList = errMsgMap
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	userCreated, err := user.SaveUser(s.DB)
	if err != nil {
		errList = formatError.FormatError(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  errList,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": userCreated,
	})
}

func (s *Server) GetUsers(c *gin.Context) {

	errList := make(map[string]string)
	user := &models.User{}
	userList, err := user.FindAllUsers(s.DB)
	if err != nil {
		errList["No_user"] = "No User Found"
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  errList,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": userList,
	})
}

func (s *Server) GetUser(c *gin.Context) {

	errList := make(map[string]string)

	userIdStr := c.Param("id")
	userId, err := strconv.ParseInt(userIdStr, 10, 32)
	if err != nil {
		errList["Invalid_request"] = "Invalid Request"
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  errList,
		})
		return
	}

	user := &models.User{}
	user, err = user.FindUserById(s.DB, uint32(userId))
	if err != nil {
		errList["No_user"] = "No User Found"
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  errList,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": user,
	})

}

func (s *Server) UpdateUser(c *gin.Context) {

	errList := make(map[string]string)

	userIdStr := c.Param("id")
	userId32, err := strconv.ParseInt(userIdStr, 10, 32)
	userId := uint32(userId32)
	if err != nil {
		errList["Invalid_request"] = "Invalid Request"
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  errList,
		})
		return
	}

	tokenId, err := auth.ExtractTokenId(c.Request)
	if err != nil {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}

	if tokenId != userId {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}

	// 开始处理请求
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		errList["Invalid_body"] = "Unable to get request"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	requestBody := make(map[string]string)
	err = json.Unmarshal(body, &requestBody)
	if err != nil {
		errList["Unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	user := &models.User{}
	dbUser, err := user.FindUserById(s.DB, userId)
	if err != nil {
		errList["User_invalid"] = "The user is does not exist"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	// When current password has content.
	if requestBody["current_password"] == "" && requestBody["new_password"] != "" {
		errList["Empty_current"] = "Please Provide current password"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	if requestBody["current_password"] != "" && requestBody["new_password"] == "" {
		errList["Empty_new"] = "Please Provide new password"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	if requestBody["current_password"] != "" && requestBody["new_password"] != "" {
		// Also check if the new password
		if len(requestBody["new_password"]) < 6 {
			errList["Invalid_password"] = "Password should be atleast 6 characters"
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"status": http.StatusUnprocessableEntity,
				"error":  errList,
			})
			return
		}

		// if they do, check that the former password is correct
		err = security.VerifyPassword(dbUser.Password, requestBody["current_password"])
		if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
			errList["Password_mismatch"] = "The password not correct"
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"status": http.StatusUnprocessableEntity,
				"error":  errList,
			})
			return
		}

		user.Password = requestBody["new_password"]
	}

	//The password fields not entered, so update only the email
	user.Username = dbUser.Username
	user.Email = requestBody["email"]

	user.Prepare()
	validateErrMsgList := user.Validate(models.USER_TYPE_UPDATE)
	if len(validateErrMsgList) > 0 {
		errList = validateErrMsgList
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	updatedUser, err := user.UpdateAUser(s.DB, userId)
	if err != nil {
		errList = formatError.FormatError(err.Error())
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": updatedUser,
	})
}

func (s *Server) UpdateAvatar(c *gin.Context) {

	errList := make(map[string]string)

	userIdStr := c.Param("id")
	userId32, err := strconv.ParseInt(userIdStr, 10, 32)
	userId := uint32(userId32)
	if err != nil {
		errList["Invalid_request"] = "Invalid Request"
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  errList,
		})
		return
	}

	tokenId, err := auth.ExtractTokenId(c.Request)
	if err != nil {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}

	if tokenId != userId {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}

	file, err := c.FormFile("avatar")
	if err != nil {
		errList["Invalid_file"] = "Invalid File"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	f, err := file.Open()
	if err != nil {
		errList["Invalid_file"] = "Invalid File"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	defer func() {
		_ = f.Close()
	}()

	size := file.Size
	// The image should not be more than 500KB
	if size > int64(512000) {
		errList["Too_large"] = "Sorry, Please upload an Image of 500KB or less"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	// 只读取512个字节就能获取文件类型
	buffer := make([]byte, 512)
	_, _ = f.Read(buffer)
	fileType := http.DetectContentType(buffer)
	if !strings.HasPrefix(fileType, "image") {
		errList["Not_Image"] = "Please Upload a valid image"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	fileBytes := bytes.NewBuffer(buffer)
	filePath := fileformat.UniqueFormat(file.Filename)
	path := "profile-photos/" + filePath
	dst, err := os.Create(path)
	if err != nil {
		errList["Upload_Failed"] = "Create file failed"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	_, err = io.Copy(dst, fileBytes)
	if err != nil {
		errList["Upload_Failed"] = "Create file failed"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	user := &models.User{}
	user.AvatarPath = filePath
	user.Prepare()
	updatedUser, err := user.UpdateAUserAvatar(s.DB, userId)
	if err != nil {
		errList["Cannot_Save"] = "Cannot Save Image, Pls try again later"
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  errList,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": updatedUser,
	})
}

func (s *Server) DeleteUser(c *gin.Context) {

	errList := make(map[string]string)

	userIdStr := c.Param("id")
	userId32, err := strconv.ParseInt(userIdStr, 10, 32)
	userId := uint32(userId32)
	if err != nil {
		errList["Invalid_request"] = "Invalid Request"
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  errList,
		})
		return
	}

	tokenId, err := auth.ExtractTokenId(c.Request)
	if err != nil {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}

	if tokenId != userId {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}

	user := &models.User{}
	_, err = user.FindUserById(s.DB, userId)
	if err != nil {
		errList["User_invalid"] = "The user is does not exist"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	err = s.DB.Transaction(func(tx *gorm.DB) error {
		_, err := user.DeleteAUser(s.DB, userId)
		if err != nil {
			return err
		}
		post := &models.Post{}
		_, err = post.DelteUserPosts(s.DB, userId)
		if err != nil {
			return err
		}
		comment := &models.Comment{}
		_, err = comment.DeleteUserComments(s.DB, userId)
		if err != nil {
			return err
		}
		like := &models.Like{}
		_, err = like.DeleteUserLikes(s.DB, userId)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		errList["Other_error"] = "Please try again later"
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": "User deleted",
	})
}
