package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/summerKK/go-code-snippet-library/champion-go/api/models"
)

func TestCreateUser(t *testing.T) {

	gin.SetMode(gin.TestMode)

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}
	samples := []struct {
		inputJSON  string
		statusCode int
		username   string
		email      string
	}{
		{
			inputJSON:  `{"username":"Pet", "email": "pet@example.com", "password": "password"}`,
			statusCode: 201,
			username:   "Pet",
			email:      "pet@example.com",
		},
		{
			inputJSON:  `{"username":"Frank", "email": "pet@example.com", "password": "password"}`,
			statusCode: 500,
		},
		{
			inputJSON:  `{"username":"Pet", "email": "grand@example.com", "password": "password"}`,
			statusCode: 500,
		},
		{
			inputJSON:  `{"username":"Kan", "email": "kanexample.com", "password": "password"}`,
			statusCode: 422,
		},
		{
			inputJSON:  `{"username": "", "email": "kan@example.com", "password": "password"}`,
			statusCode: 422,
		},
		{
			inputJSON:  `{"username": "Kan", "email": "", "password": "password"}`,
			statusCode: 422,
		},
		{
			inputJSON:  `{"username": "Kan", "email": "kan@example.com", "password": ""}`,
			statusCode: 422,
		},
	}

	r := gin.Default()
	r.POST("/users", server.CreateUser)
	for _, sample := range samples {
		req, err := http.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(sample.inputJSON))
		if err != nil {
			t.Error(err)
			return
		}

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		responseInterface := make(map[string]interface{})
		err = json.Unmarshal(rr.Body.Bytes(), &responseInterface)
		if err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, rr.Code, sample.statusCode)

		if rr.Code == http.StatusCreated {
			responseInfo := responseInterface["response"].(map[string]interface{})
			assert.Equal(t, responseInfo["username"], sample.username)
			assert.Equal(t, responseInfo["email"], sample.email)
		}

		if sample.statusCode == 422 || sample.statusCode == 500 {
			responseMap := responseInterface["error"].(map[string]interface{})

			if responseMap["Taken_email"] != nil {
				assert.Equal(t, responseMap["Taken_email"], "Email Already Taken")
			}
			if responseMap["Taken_username"] != nil {
				assert.Equal(t, responseMap["Taken_username"], "Username Already Taken")
			}
			if responseMap["Invalid_email"] != nil {
				assert.Equal(t, responseMap["Invalid_email"], "Invalid Email")
			}
			if responseMap["Required_username"] != nil {
				assert.Equal(t, responseMap["Required_username"], "Required Username")
			}
			if responseMap["Required_email"] != nil {
				assert.Equal(t, responseMap["Required_email"], "Required Email")
			}
			if responseMap["Required_password"] != nil {
				assert.Equal(t, responseMap["Required_password"], "Required Password")
			}
		}
	}
}

func TestGetUsers(t *testing.T) {

	gin.SetMode(gin.TestMode)

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	userList, err := seedUsers()
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	r.GET("/users", server.GetUsers)

	req, err := http.NewRequest(http.MethodGet, "/users", nil)
	if err != nil {
		t.Error(err)
		return
	}

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("unexcept http status code %d\n", rr.Code)
		return
	}

	responseInterface := make(map[string]interface{})
	err = json.Unmarshal(rr.Body.Bytes(), &responseInterface)
	if err != nil {
		t.Error(err)
		return
	}

	users := responseInterface["response"].([]interface{})

	assert.Equal(t, len(userList), len(users))
}

func TestGetUser(t *testing.T) {

	gin.SetMode(gin.TestMode)

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	user, err := seedOneUser()
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	r.GET("/users/:id", server.GetUser)

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/users/%d", user.ID), nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("unexcept http status code %d\n", rr.Code)
		return
	}

	responseInterface := make(map[string]interface{})
	err = json.Unmarshal(rr.Body.Bytes(), &responseInterface)
	if err != nil {
		t.Error(err)
		return
	}
	responseInfo := responseInterface["response"].(map[string]interface{})

	assert.Equal(t, user.ID, uint32(responseInfo["id"].(float64)))
	assert.Equal(t, user.Username, responseInfo["username"])
	assert.Equal(t, user.Email, responseInfo["email"])
}

func TestUpdateUser(t *testing.T) {

	gin.SetMode(gin.TestMode)

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	user, err := seedOneUser()
	if err != nil {
		log.Fatal(err)
	}

	userUpdate := &models.User{
		ID:       user.ID,
		Username: "Summer",
		Email:    "summer@qq.com",
		Password: "password",
	}

	updatedUser, err := userUpdate.UpdateAUser(server.DB, user.ID)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, updatedUser.ID, userUpdate.ID)
	assert.Equal(t, updatedUser.Username, userUpdate.Username)
	assert.Equal(t, updatedUser.Email, userUpdate.Email)
}
