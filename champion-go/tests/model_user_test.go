package tests

import (
	"log"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/jinzhu/gorm"
	"github.com/summerKK/go-code-snippet-library/champion-go/api/models"
	"github.com/summerKK/go-code-snippet-library/champion-go/api/security"
)

func TestFindAllUsers(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}
	sUsers, err := seedUsers()
	if err != nil {
		log.Fatal(err)
	}

	users, err := userInstance.FindAllUsers(server.DB)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, len(sUsers), len(users))
}

func TestSaveUser(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	newUser := &models.User{
		ID:       1,
		Email:    "test@example.com",
		Username: "test",
		Password: "password",
	}

	savedUser, err := newUser.SaveUser(server.DB)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, newUser.ID, savedUser.ID)
	assert.Equal(t, newUser.Email, savedUser.Email)
	assert.Equal(t, newUser.Username, savedUser.Username)
	assert.Equal(t, newUser.Password, savedUser.Password)
}

func TestFindUserById(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	user, err := seedOneUser()
	if err != nil {
		log.Fatal(err)
	}

	foundUser, err := userInstance.FindUserById(server.DB, user.ID)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, user.ID, foundUser.ID)
	assert.Equal(t, user.Email, foundUser.Email)
	assert.Equal(t, user.Username, foundUser.Username)
	assert.Equal(t, user.Password, foundUser.Password)
}

func TestUpdatAUser(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	user, err := seedOneUser()
	if err != nil {
		log.Fatal(err)
	}

	userUpdate := &models.User{
		ID:       1,
		Username: "modiUpdate",
		Email:    "modiupdate@example.com",
		Password: "password",
	}

	updatedAUser, err := userUpdate.UpdateAUser(server.DB, user.ID)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, userUpdate.ID, updatedAUser.ID)
	assert.Equal(t, userUpdate.Email, updatedAUser.Email)
	assert.Equal(t, userUpdate.Username, updatedAUser.Username)
	assert.Equal(t, userUpdate.Password, updatedAUser.Password)
}

func TestUpdateAUserAvatar(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	user, err := seedOneUser()
	if err != nil {
		log.Fatal(err)
	}

	userUpdate := &models.User{
		AvatarPath: "hello,world",
	}

	updatedUser, err := userUpdate.UpdateAUserAvatar(server.DB, user.ID)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, userUpdate.AvatarPath, updatedUser.AvatarPath)
}

func TestDeleteAUser(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	user, err := seedOneUser()
	if err != nil {
		log.Fatal(err)
	}

	_, err = user.DeleteAUser(server.DB, user.ID)
	if err != nil {
		t.Error(err)
		return
	}

	_, err = userInstance.FindUserById(server.DB, user.ID)

	if !gorm.IsRecordNotFoundError(err) {
		t.Error(err)
	}
}

func TestUpdateAUserPassword(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	user, err := seedOneUser()
	if err != nil {
		log.Fatal(err)
	}

	passwordStr := "hello,world"
	userUpdate := &models.User{
		Password: passwordStr,
		Email:    user.Email,
	}

	updatedUser, err := userUpdate.UpdateAUserPassword(server.DB)
	if err != nil {
		t.Error(err)
		return
	}

	err = security.VerifyPassword(updatedUser.Password, passwordStr)
	if err != nil {
		t.Error(err)
		return
	}

}
