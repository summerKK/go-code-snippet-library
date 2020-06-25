package models

import (
	"errors"
	"fmt"
	"html"
	"os"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"github.com/jinzhu/gorm"
	"github.com/summerKK/go-code-snippet-library/champion-go/api/security"
)

type User struct {
	ID         uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Username   string    `gorm:"size:255;not null;unique" json:"username"`
	Email      string    `gorm:"size:100;not null;unique" json:"email"`
	Password   string    `gorm:"size:100;not null;" json:"password"`
	AvatarPath string    `gorm:"size:255;null;" json:"avatar_path"`
	CreatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

type UserType int

const (
	USER_TYPE_UPDATE UserType = iota
	USER_TYPE_LOGIN
	USER_TYPE_FORGOT_PASSWORD
	USER_TYPE_DEFAULT
)

// 钩子函数.保存数据的时候自动触发
func (u *User) BeforeSave() error {

	hashedPassword, err := security.Hash(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)

	return nil
}

func (u *User) AfterFind() error {

	if u.AvatarPath != "" {
		u.AvatarPath = os.Getenv("DO_SPACES_URL") + u.AvatarPath
	}

	return nil
}

func (u *User) Prepare() {

	u.Username = html.EscapeString(strings.TrimSpace(u.Username))
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
}

func (u *User) Validate(action UserType) map[string]string {

	var err error
	errMsgList := make(map[string]string)

	switch action {
	case USER_TYPE_UPDATE:
		if u.Email == "" {
			err = errors.New("Required Email")
			errMsgList["Required_email"] = err.Error()
		}
		if u.Email != "" {
			if err = checkmail.ValidateFormat(u.Email); err != nil {
				err = errors.New("Invalid Email")
				errMsgList["Invalid_email"] = err.Error()
			}
		}

	case USER_TYPE_LOGIN:
		if u.Password == "" {
			err = errors.New("Required Password")
			errMsgList["Required_password"] = err.Error()
		}
		if u.Email == "" {
			err = errors.New("Required Email")
			errMsgList["Required_email"] = err.Error()
		}
		if u.Email != "" {
			if err = checkmail.ValidateFormat(u.Email); err != nil {
				err = errors.New("Invalid Email")
				errMsgList["Invalid_email"] = err.Error()
			}
		}

	case USER_TYPE_FORGOT_PASSWORD:
		if u.Email == "" {
			err = errors.New("Required Email")
			errMsgList["Required_email"] = err.Error()
		}
		if u.Email != "" {
			if err = checkmail.ValidateFormat(u.Email); err != nil {
				err = errors.New("Invalid Email")
				errMsgList["Invalid_email"] = err.Error()
			}
		}

	default:
		if u.Username == "" {
			err = errors.New("Required Username")
			errMsgList["Required_username"] = err.Error()
		}
		if u.Password == "" {
			err = errors.New("Required Password")
			errMsgList["Required_password"] = err.Error()
		}
		if u.Password != "" && len(u.Password) < 6 {
			err = errors.New("Password should be atleast 6 characters")
			errMsgList["Invalid_password"] = err.Error()
		}
		if u.Email == "" {
			err = errors.New("Required Email")
			errMsgList["Required_email"] = err.Error()

		}
		if u.Email != "" {
			if err = checkmail.ValidateFormat(u.Email); err != nil {
				err = errors.New("Invalid Email")
				errMsgList["Invalid_email"] = err.Error()
			}
		}
	}

	return errMsgList
}

func (u *User) SaveUser(db *gorm.DB) (*User, error) {

	var err error
	err = db.Debug().Create(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) FindAllUsers(db *gorm.DB) (users []*User, err error) {
	err = db.Debug().Model(&User{}).Limit(100).Find(&users).Error

	return
}

func (u *User) FindUserById(db *gorm.DB, uid uint32) (user *User, err error) {

	user = &User{}
	err = db.Debug().Model(User{}).Where("id = ?", uid).Take(user).Error
	if err != nil {
		return
	}

	return
}

func (u *User) FindUserByEmail(db *gorm.DB, email string) (user *User, err error) {

	user = &User{}
	err = db.Debug().Model(User{}).Where("email = ?", email).Take(user).Error
	if err != nil {
		return
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, fmt.Errorf("User not found(user_id:%d)", email)
	}

	return
}

func (u *User) UpdateAUser(db *gorm.DB, uid uint32) (user *User, err error) {

	user = &User{}
	err = db.Debug().Model(User{}).Where("id = ?", uid).Updates(
		map[string]interface{}{
			"password":  u.Password,
			"username":  u.Username,
			"email":     u.Email,
			"update_at": time.Now(),
		},
	).Error

	if err != nil {
		return
	}
	err = db.Debug().Model(User{}).Where("id = ?", uid).Take(user).Error

	return
}

func (u *User) UpdateAUserAvatar(db *gorm.DB, uid uint32) (user *User, err error) {

	user = &User{}
	db = db.Debug().Model(User{}).Where("id = ?", uid).Take(&User{}).Updates(
		map[string]interface{}{
			"avatar_path": u.AvatarPath,
			"update_at":   time.Now(),
		},
	)

	if db.Error != nil {
		err = db.Error
		return
	}
	err = db.Debug().Model(User{}).Where("id = ?", uid).Take(user).Error

	return
}

func (u *User) DeleteAUser(db *gorm.DB, uid uint32) (rowAffected int64, err error) {

	db = db.Debug().Model(User{}).Where("id = ?", uid).Take(&User{}).Delete(&User{})
	if db.Error == nil {
		rowAffected = db.RowsAffected
	}

	return
}

func (u *User) UpdateAUserPassword(db *gorm.DB) (user *User, err error) {

	user = &User{}
	// hash the password
	err = u.BeforeSave()
	if err != nil {
		return
	}

	err = db.Debug().Model(User{}).Where("email = ?", u.Email).Updates(
		map[string]interface{}{
			"password":  u.Password,
			"update_at": time.Now(),
		},
	).Error

	if db.Error != nil {
		err = db.Error
		return
	}

	err = db.Debug().Model(User{}).Where("email", u.Email).Take(user).Error

	return
}
