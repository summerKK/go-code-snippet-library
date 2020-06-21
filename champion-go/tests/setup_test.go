package tests

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/summerKK/go-code-snippet-library/champion-go/api/controllers"
	"github.com/summerKK/go-code-snippet-library/champion-go/api/models"
)

var server = &controllers.Server{}
var userInstance = &models.User{}
var postInstance = &models.Post{}
var likeInstance = &models.Like{}
var commentInstance = &models.Comment{}

func TestMain(m *testing.M) {
	var err error
	err = godotenv.Load(os.ExpandEnv("../.env"))
	if err != nil {
		log.Fatal(err)
	}

	Database()
	os.Exit(m.Run())
}

func Database() {
	var err error
	TestDbDriver := os.Getenv("TEST_DB_DRIVER")
	DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", os.Getenv("TEST_DB_USER"), os.Getenv("TEST_DB_PASSWORD"), os.Getenv("TEST_DB_HOST"), os.Getenv("TEST_DB_PORT"), os.Getenv("TEST_DB_NAME"))
	server.DB, err = gorm.Open(TestDbDriver, DBURL)
	if err != nil {
		fmt.Printf("Cannot connect to %s database\n", TestDbDriver)
		log.Fatal("This is the error:", err)
	} else {
		fmt.Printf("We are connected to the %s database\n", TestDbDriver)
	}
}

func refreshUserTable() (err error) {
	err = server.DB.DropTableIfExists(models.User{}).Error
	if err != nil {
		return
	}

	err = server.DB.AutoMigrate(models.User{}).Error
	if err != nil {
		return
	}

	log.Printf("Successfully refreshed table")
	return nil
}

func seedOneUser() (user *models.User, err error) {
	user = &models.User{
		Username: "Pet",
		Email:    "pet@example.com",
		Password: "password",
	}
	err = server.DB.Model(models.User{}).Create(&user).Error
	if err != nil {
		return
	}

	return user, nil
}

func seedUsers() (users []*models.User, err error) {
	users = []*models.User{
		{
			Username: "Steven",
			Email:    "steven@example.com",
			Password: "password",
		},
		{
			Username: "Kenny",
			Email:    "kenny@example.com",
			Password: "password",
		},
	}

	for _, user := range users {
		err = server.DB.Model(models.User{}).Create(user).Error
		if err != nil {
			return
		}
	}

	return
}

func refreshUserAndPostTable() error {

	err := server.DB.DropTableIfExists(&models.User{}, &models.Post{}).Error
	if err != nil {
		return err
	}
	err = server.DB.AutoMigrate(&models.User{}, &models.Post{}).Error
	if err != nil {
		return err
	}
	log.Printf("Successfully refreshed tables")
	return nil
}

func seedOneUserAndOnePost() (user *models.User, post *models.Post, err error) {
	user = &models.User{
		Username: "Sam",
		Email:    "sam@example.com",
		Password: "password",
	}
	err = server.DB.Model(models.User{}).Create(&user).Error
	if err != nil {
		return
	}

	post = &models.Post{
		Title:    "This is the title sam",
		Content:  "This is the content sam",
		AuthorID: user.ID,
	}
	err = server.DB.Model(models.Post{}).Create(&post).Error
	if err != nil {
		return
	}

	return
}

func seedUsersAndPosts() (users []*models.User, posts []*models.Post, err error) {

	users = []*models.User{
		{
			Username: "Steven",
			Email:    "steven@example.com",
			Password: "password",
		},
		{
			Username: "Magu",
			Email:    "magu@example.com",
			Password: "password",
		},
	}
	posts = []*models.Post{
		{
			Title:   "Title 1",
			Content: "Hello world 1",
		},
		{
			Title:   "Title 2",
			Content: "Hello world 2",
		},
	}

	for i := range users {
		err = server.DB.Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			return
		}
		posts[i].AuthorID = users[i].ID

		err = server.DB.Model(&models.Post{}).Create(&posts[i]).Error
		if err != nil {
			return
		}
	}

	return
}

func refreshUserPostAndLikeTable() (err error) {
	err = server.DB.DropTableIfExists(&models.User{}, &models.Post{}, &models.Like{}).Error
	if err != nil {
		return
	}
	err = server.DB.AutoMigrate(&models.User{}, &models.Post{}, &models.Like{}).Error
	if err != nil {
		return
	}
	log.Printf("Successfully refreshed user, post and like tables")

	return
}

func seedUsersPostsAndLikes() (post *models.Post, users []*models.User, likes []*models.Like, err error) {
	// The idea here is: two users can like one post
	users = []*models.User{
		{
			Username: "Steven",
			Email:    "steven@example.com",
			Password: "password",
		},
		{
			Username: "Magu",
			Email:    "magu@example.com",
			Password: "password",
		},
	}
	post = &models.Post{
		Title:   "This is the title",
		Content: "This is the content",
	}
	err = server.DB.Model(&models.Post{}).Create(&post).Error
	if err != nil {
		log.Fatalf("cannot seed post table: %v", err)
	}
	likes = []*models.Like{
		{
			UserID: 1,
			PostID: post.ID,
		},
		{
			UserID: 2,
			PostID: post.ID,
		},
	}
	for i, _ := range users {
		err = server.DB.Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			return
		}
		err = server.DB.Model(&models.Like{}).Create(&likes[i]).Error
		if err != nil {
			return
		}
	}

	return
}

func refreshUserPostAndCommentTable() (err error) {
	err = server.DB.DropTableIfExists(&models.User{}, &models.Post{}, &models.Comment{}).Error
	if err != nil {
		return
	}
	err = server.DB.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{}).Error
	if err != nil {
		return
	}
	log.Printf("Successfully refreshed user, post and comment tables")

	return
}

func seedUsersPostsAndComments() (post *models.Post, users []*models.User, comments []*models.Comment, err error) {
	// The idea here is: two users can comment one post
	users = []*models.User{
		{
			Username: "Steven",
			Email:    "steven@example.com",
			Password: "password",
		},
		{
			Username: "Magu",
			Email:    "magu@example.com",
			Password: "password",
		},
	}
	post = &models.Post{
		Title:   "This is the title",
		Content: "This is the content",
	}
	err = server.DB.Model(&models.Post{}).Create(&post).Error
	if err != nil {
		log.Fatalf("cannot seed post table: %v", err)
	}
	comments = []*models.Comment{
		{
			Body:   "user 1 made this comment",
			UserID: 1,
			PostID: post.ID,
		},
		{
			Body:   "user 2 made this comment",
			UserID: 2,
			PostID: post.ID,
		},
	}
	for i, _ := range users {
		err = server.DB.Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			return
		}
		err = server.DB.Model(&models.Like{}).Create(&comments[i]).Error
		if err != nil {
			return
		}
	}

	return
}

func refreshUserAndResetPasswordTable() (err error) {
	err = server.DB.DropTableIfExists(&models.User{}, &models.ResetPassword{}).Error
	if err != nil {
		return
	}
	err = server.DB.AutoMigrate(&models.User{}, &models.ResetPassword{}).Error
	if err != nil {
		return
	}
	log.Printf("Successfully refreshed user and resetpassword tables")

	return
}

// Seed the reset password table with the token
func seedResetPassword() (resetPassword *models.ResetPassword, err error) {

	resetPassword = &models.ResetPassword{
		Token: "awesometoken",
		Email: "pet@example.com",
	}
	err = server.DB.Model(&models.ResetPassword{}).Create(&resetPassword).Error
	if err != nil {
		return
	}

	return
}
