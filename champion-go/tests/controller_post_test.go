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
)

func TestCreatePost(t *testing.T) {

	gin.SetMode(gin.TestMode)

	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatal(err)
	}

	user, err := seedOneUser()
	if err != nil {
		log.Fatal(err)
	}

	tokenInterface, err := server.SignIn(user.Email, "password")
	if err != nil {
		log.Fatal(err)
	}
	token := tokenInterface["token"]
	tokenString := fmt.Sprintf("Bearer %v", token)

	// Note that the author id is obtained from the token, so we dont pass it
	samples := []struct {
		inputJSON  string
		statusCode int
		title      string
		content    string
		tokenGiven string
	}{
		{
			inputJSON:  `{"title":"The title", "content": "the content"}`,
			statusCode: 201,
			tokenGiven: tokenString,
			title:      "The title",
			content:    "the content",
		},
		{
			// When the post title already exist
			inputJSON:  `{"title":"The title", "content": "the content"}`,
			statusCode: 500,
			tokenGiven: tokenString,
		},
		{
			// When no token is passed
			inputJSON:  `{"title":"When no token is passed", "content": "the content"}`,
			statusCode: 401,
			tokenGiven: "",
		},
		{
			// When incorrect token is passed
			inputJSON:  `{"title":"When incorrect token is passed", "content": "the content"}`,
			statusCode: 401,
			tokenGiven: "This is an incorrect token",
		},
		{
			inputJSON:  `{"title": "", "content": "The content"}`,
			statusCode: 422,
			tokenGiven: tokenString,
		},
		{
			inputJSON:  `{"title": "This is a title", "content": ""}`,
			statusCode: 422,
			tokenGiven: tokenString,
		},
	}

	r := gin.Default()
	r.POST("/posts", server.CreatePost)
	for _, sample := range samples {
		req, err := http.NewRequest(http.MethodPost, "/posts", bytes.NewBufferString(sample.inputJSON))
		if err != nil {
			t.Error(err)
			return
		}
		req.Header.Set("Authorization", sample.tokenGiven)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		responseInterface := make(map[string]interface{})
		assert.Equal(t, rr.Code, sample.statusCode)

		err = json.Unmarshal(rr.Body.Bytes(), &responseInterface)
		if err != nil {
			t.Error(err)
			return
		}

		if rr.Code == http.StatusCreated {
			responseInfo := responseInterface["response"].(map[string]interface{})
			assert.Equal(t, sample.content, responseInfo["content"])
			assert.Equal(t, sample.title, responseInfo["title"])
		}

		if sample.statusCode == 401 || sample.statusCode == 422 || sample.statusCode == 500 {
			responseMap := responseInterface["error"].(map[string]interface{})

			if responseMap["Unauthorized"] != nil {
				assert.Equal(t, responseMap["Unauthorized"], "Unauthorized")
			}
			if responseMap["Taken_title"] != nil {
				assert.Equal(t, responseMap["Taken_title"], "Title Already Taken")
			}
			if responseMap["Required_title"] != nil {
				assert.Equal(t, responseMap["Required_title"], "Required Title")
			}
			if responseMap["Required_content"] != nil {
				assert.Equal(t, responseMap["Required_content"], "Required Content")
			}
		}
	}
}

func TestGetPosts(t *testing.T) {

	gin.SetMode(gin.TestMode)

	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatal(err)
	}

	_, posts, err := seedUsersAndPosts()
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	r.GET("/posts", server.GetPosts)

	req, err := http.NewRequest(http.MethodGet, "/posts", nil)
	if err != nil {
		t.Error(err)
		return
	}

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusOK)

	responseInterface := make(map[string]interface{})
	err = json.Unmarshal(rr.Body.Bytes(), &responseInterface)
	if err != nil {
		t.Error(err)
		return
	}

	thePosts := responseInterface["response"].([]interface{})
	assert.Equal(t, len(thePosts), len(posts))
}

func TestGetPost(t *testing.T) {

	gin.SetMode(gin.TestMode)

	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatal(err)
	}

	_, post, err := seedOneUserAndOnePost()
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	r.GET("/posts/:id", server.GetPost)

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/posts/%d", post.ID), nil)
	if err != nil {
		t.Error(err)
		return
	}

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusOK)

	responseInterface := make(map[string]interface{})
	err = json.Unmarshal(rr.Body.Bytes(), &responseInterface)
	if err != nil {
		t.Error(err)
		return
	}

	getPost := responseInterface["response"].(map[string]interface{})
	assert.Equal(t, getPost["title"], post.Title)
	assert.Equal(t, getPost["content"], post.Content)
	assert.Equal(t, uint32(getPost["author_id"].(float64)), post.AuthorID)
}

func TestUpdatePost(t *testing.T) {

	gin.SetMode(gin.TestMode)

	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatal(err)
	}

	users, posts, err := seedUsersAndPosts()
	if err != nil {
		log.Fatal(err)
	}

	tokenInterface0, err := server.SignIn(users[0].Email, "password")
	token0 := tokenInterface0["token"]
	tokenString0 := fmt.Sprintf("Bearer %s", token0)

	tokenInterface1, err := server.SignIn(users[1].Email, "password")
	token1 := tokenInterface1["token"]
	tokenString1 := fmt.Sprintf("Bearer %s", token1)

	samples := []struct {
		id         uint64
		updateJSON string
		statusCode int
		title      string
		content    string
		tokenGiven string
	}{
		{
			// Convert int64 to int first before converting to string
			id:         posts[0].ID,
			updateJSON: `{"title":"The updated posts", "content": "This is the updated content"}`,
			statusCode: 200,
			title:      "The updated posts",
			content:    "This is the updated content",
			tokenGiven: tokenString0,
		},
		{
			// When no token0 is provided
			id:         posts[0].ID,
			updateJSON: `{"title":"This is still another title", "content": "This is the updated content"}`,
			tokenGiven: "",
			statusCode: 401,
		},
		{
			// When incorrect token0 is provided
			id:         posts[0].ID,
			updateJSON: `{"title":"This is still another title", "content": "This is the updated content"}`,
			tokenGiven: "this is an incorrect token",
			statusCode: 401,
		},
		{
			//Note: "Title 2" belongs to posts 2, and title must be unique
			id:         posts[1].ID,
			updateJSON: `{"title":"The updated posts", "content": "This is the updated content"}`,
			statusCode: 500,
			tokenGiven: tokenString1,
		},
		{
			// When title is not given
			id:         posts[0].ID,
			updateJSON: `{"title":"", "content": "This is the updated content"}`,
			statusCode: 422,
			tokenGiven: tokenString0,
		},
		{
			// When content is not given
			id:         posts[0].ID,
			updateJSON: `{"title":"Awesome title", "content": ""}`,
			statusCode: 422,
			tokenGiven: tokenString0,
		},
		{
			// When invalid posts id is given
			id:         uint64(0),
			statusCode: 401,
		},
	}

	r := gin.Default()
	r.PUT("/posts/:id", server.UpdatePost)
	for _, sample := range samples {
		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/posts/%d", sample.id), bytes.NewBufferString(sample.updateJSON))
		if err != nil {
			t.Error(err)
			return
		}
		req.Header.Set("Authorization", sample.tokenGiven)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		assert.Equal(t, sample.statusCode, rr.Code)
	}
}
