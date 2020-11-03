package model

import "github.com/jinzhu/gorm"

type User struct {
	Model
	Name          string `json:"name"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	IsAdmin       int8   `json:"is_admin"`
	Preferences   string `json:"preferences"`
	RememberToken string `json:"remember_token"`
}

func (u *User) TableName() string {
	return "users"
}

func (u *User) GetUserByEmail(db *gorm.DB) (*User, error) {
	var user User
	err := db.Where("email = ?", u.Email).First(&user).Error

	return &user, err
}
