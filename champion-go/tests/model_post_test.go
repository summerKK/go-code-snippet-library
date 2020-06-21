package tests

import (
	"log"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/summerKK/go-code-snippet-library/champion-go/api/models"
)

func TestFindAllPosts(t *testing.T) {
	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatal(err)
	}

	_, posts, err := seedUsersAndPosts()
	if err != nil {
		log.Fatal(err)
	}

	allPosts, err := postInstance.FindAllPosts(server.DB)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, len(posts), len(allPosts))
}

func TestSavePost(t *testing.T) {
	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatal(err)
	}

	user, err := seedOneUser()
	if err != nil {
		log.Fatal(err)
	}

	post := &models.Post{
		ID:       1,
		Title:    "This is the title sam",
		Content:  "This is the content sam",
		AuthorID: user.ID,
	}

	savedPost, err := post.SavePost(server.DB)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, post.ID, savedPost.ID)
	assert.Equal(t, post.Title, savedPost.Title)
	assert.Equal(t, post.Content, savedPost.Content)
	assert.Equal(t, post.AuthorID, savedPost.AuthorID)
}

func TestFindPostById(t *testing.T) {
	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatal(err)
	}

	_, post, err := seedOneUserAndOnePost()
	if err != nil {
		log.Fatal(err)
	}

	foundPost, err := postInstance.FindPostById(server.DB, post.ID)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, post.ID, foundPost.ID)
	assert.Equal(t, post.Title, foundPost.Title)
	assert.Equal(t, post.Content, foundPost.Content)
	assert.Equal(t, post.AuthorID, foundPost.AuthorID)
}

func TestUpdateAPost(t *testing.T) {
	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatal(err)
	}

	_, post, err := seedOneUserAndOnePost()
	if err != nil {
		log.Fatal(err)
	}

	updatePost := &models.Post{
		ID:       1,
		Title:    "modiUpdate",
		Content:  "modiupdate@example.com",
		AuthorID: post.AuthorID,
	}

	updatedPost, err := updatePost.UpdateAPost(server.DB)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, updatePost.ID, updatedPost.ID)
	assert.Equal(t, updatePost.Title, updatedPost.Title)
	assert.Equal(t, updatePost.Content, updatedPost.Content)
	assert.Equal(t, updatePost.AuthorID, updatedPost.AuthorID)
}

func TestDeleteAPost(t *testing.T) {
	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatal(err)
	}

	_, post, err := seedOneUserAndOnePost()
	if err != nil {
		log.Fatal(err)
	}

	rowAffected, err := postInstance.DeleteAPost(server.DB, post.ID, post.AuthorID)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, int64(1), rowAffected)
}

func TestFindUserPosts(t *testing.T) {
	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatal(err)
	}

	users, _, err := seedUsersAndPosts()
	if err != nil {
		log.Fatal(err)
	}

	userPosts, err := postInstance.FindUserPosts(server.DB, users[0].ID)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, len(userPosts), 1)
}

func TestDelteUserPosts(t *testing.T) {
	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatal(err)
	}

	users, _, err := seedUsersAndPosts()
	if err != nil {
		log.Fatal(err)
	}

	rowAffected, err := postInstance.DelteUserPosts(server.DB, users[0].ID)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, rowAffected, int64(1))
}
