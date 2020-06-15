package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/summerKK/go-code-snippet-library/champion-go/api/auth"
	"github.com/summerKK/go-code-snippet-library/champion-go/api/models"
	"github.com/summerKK/go-code-snippet-library/champion-go/api/security"
	"github.com/summerKK/go-code-snippet-library/champion-go/api/utils/formatError"
)

func (s *Server) Login(c *gin.Context) {

	errList := make(map[string]string)

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		errList["Invalid_request"] = "Invalid Request"
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  errList,
		})
		return
	}

	user := &models.User{}
	err = json.Unmarshal(body, user)
	if err != nil {
		errList["Unmarshal_error"] = "Unmarshal error"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	user.Prepare()
	validateMsgList := user.Validate(models.USER_TYPE_LOGIN)
	if len(validateMsgList) > 0 {
		errList = validateMsgList
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	userInfo, err := s.signIn(user.Email, user.Password)
	if err != nil {
		formattedError := formatError.FormatError(err.Error())
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  formattedError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": userInfo,
	})
}

func (s *Server) signIn(email, password string) (map[string]interface{}, error) {

	userData := make(map[string]interface{})

	user := &models.User{}

	user, err := user.FindUserByEmail(s.DB, email)
	if err != nil {
		fmt.Println("this is the error getting the user: ", err)
		return nil, err
	}

	err = security.VerifyPassword(user.Password, password)
	if err != nil {
		fmt.Println("this is the error hashing the password: ", err)
		return nil, err
	}

	token, err := auth.CreateToken(user.ID)
	if err != nil {
		fmt.Println("this is the error creating the token: ", err)
		return nil, err
	}

	userData["token"] = token
	userData["id"] = user.ID
	userData["email"] = user.Email
	userData["avatar_path"] = user.AvatarPath
	userData["username"] = user.Username

	return userData, nil
}
