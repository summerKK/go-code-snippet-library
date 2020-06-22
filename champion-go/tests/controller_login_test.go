package tests

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
)

func TestSignIn(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	user, err := seedOneUser()
	if err != nil {
		log.Fatal(err)
	}

	samples := []struct {
		email        string
		password     string
		errorMessage string
	}{
		{
			email:        user.Email,
			password:     "password", // Note the password has to be this, not the hashed one from the database
			errorMessage: "",
		},
		{
			email:        user.Email,
			password:     "Wrong password",
			errorMessage: "crypto/bcrypt: hashedPassword is not the hash of the given password",
		},
		{
			email:        "Wrong email",
			password:     "password",
			errorMessage: "record not found",
		},
	}

	for _, sample := range samples {
		loginDetails, err := server.SignIn(sample.email, sample.password)
		if err != nil {
			assert.Equal(t, err, errors.New(sample.errorMessage))
		} else {
			assert.NotEqual(t, loginDetails, "")
		}
	}
}

func TestLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	user, err := seedOneUser()
	if err != nil {
		log.Fatal(err)
	}

	samples := []struct {
		inputJSON  string
		statusCode int
		username   string
		email      string
		password   string
	}{
		{
			inputJSON:  `{"email": "pet@example.com", "password": "password"}`,
			statusCode: 200,
			username:   user.Username,
			email:      user.Email,
		},
		{
			inputJSON:  `{"email": "pet@example.com", "password": "wrong password"}`,
			statusCode: 422,
		},
		{
			// this record does not exist
			inputJSON:  `{"email": "frank@example.com", "password": "password"}`,
			statusCode: 422,
		},
		{
			inputJSON:  `{"email": "kanexample.com", "password": "password"}`,
			statusCode: 422,
		},
		{
			inputJSON:  `{"email": "", "password": "password"}`,
			statusCode: 422,
		},
		{
			inputJSON:  `{"email": "kan@example.com", "password": ""}`,
			statusCode: 422,
		},
	}

	r := gin.Default()
	r.POST("/login", server.Login)
	for _, sample := range samples {
		req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(sample.inputJSON))
		if err != nil {
			t.Error(err)
			return
		}
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		responseInterface := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseInterface)
		if err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, sample.statusCode, rr.Code)

		if sample.statusCode == 200 {
			responseMap := responseInterface["response"].(map[string]interface{})
			assert.Equal(t, responseMap["username"], sample.username)
			assert.Equal(t, responseMap["email"], sample.email)
		}

		if sample.statusCode == 422 {
			responseMap := responseInterface["error"].(map[string]interface{})

			if responseMap["Incorrect_password"] != nil {
				assert.Equal(t, responseMap["Incorrect_password"], "Incorrect Password")
			}
			if responseMap["Incorrect_details"] != nil {
				assert.Equal(t, responseMap["Incorrect_details"], "Incorrect Details")
			}
			if responseMap["Invalid_email"] != nil {
				assert.Equal(t, responseMap["Invalid_email"], "Invalid Email")
			}
			if responseMap["Required_password"] != nil {
				assert.Equal(t, responseMap["Required_password"], "Required Password")
			}
		}
	}
}
