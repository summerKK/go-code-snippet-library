package models

import (
	"errors"
	"fmt"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Comment struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	UserID    uint32    `gorm:"not null" json:"user_id"`
	PostID    uint64    `gorm:"not null" json:"post_id"`
	Body      string    `gorm:"text;not null;" json:"body"`
	User      User      `json:"user"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

type CommentType int

const (
	COMMENT_TYPE_UPDATE = iota
	COMMENT_TYPE_DEFAULT
)

func (c *Comment) Prepare() {

	c.Body = html.EscapeString(strings.TrimSpace(c.Body))
	c.User = User{}
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
}

func (c *Comment) Validate(action CommentType) map[string]string {

	var errorMessages = make(map[string]string)
	var err error

	switch action {
	case COMMENT_TYPE_UPDATE:
		if c.Body == "" {
			err = errors.New("Required Comment")
			errorMessages["Required_body"] = err.Error()
		}

	default:
		if c.Body == "" {
			err = errors.New("Required Comment")
			errorMessages["Required_body"] = err.Error()
		}
	}

	return errorMessages
}

func (c *Comment) SaveComment(db *gorm.DB) (comment *Comment, err error) {

	comment = &Comment{}
	err = db.Debug().Create(comment).Error
	if err != nil {
		return
	}
	if c.ID > 0 {
		err = db.Debug().Model(User{}).Where("id = ?", c.UserID).Take(&c.User).Error
		if err != nil {
			return
		}
	}

	return
}

func (c *Comment) GetComments(db *gorm.DB, pid uint64) (commentList []*Comment, err error) {

	err = db.Debug().Model(&Comment{}).Where("post_id = ?", pid).Order("created_at desc").Find(&commentList).Error
	if err != nil {
		return
	}
	if len(commentList) > 0 {
		for _, comment := range commentList {
			err = db.Debug().Model(User{}).Where("id = ?", comment.UserID).Take(comment.User).Error
			if err != nil {
				return
			}
		}
	}

	return
}

func (c *Comment) UpdateAComment(db *gorm.DB) (comment *Comment, err error) {

	err = db.Debug().Model(Comment{}).Where("id = ?", c.ID).Updates(Comment{Body: c.Body, UpdatedAt: time.Now()}).Error
	if err != nil {
		return
	}

	fmt.Println("this is the comment body: ", c.Body)
	if c.ID >= 0 {
		err = db.Debug().Model(User{}).Where("id = ?", c.UserID).Take(&c.User).Error
		if err != nil {
			return &Comment{}, err
		}
	}
	comment = c

	return
}

func (c *Comment) DeleteAComment(db *gorm.DB) (rowAffected int64, err error) {

	db = db.Debug().Model(Comment{}).Where("id = ?", c.ID).Delete(&Comment{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(err) {
			err = errors.New("record not found")
		}

		return
	}

	rowAffected = db.RowsAffected

	return
}

func (c *Comment) DeleteUserComments(db *gorm.DB, uid uint32) (rowAffected int64, err error) {

	db = db.Debug().Model(Comment{}).Where("user_id = ?", uid).Delete(Comment{})
	if db.Error != nil {
		if gorm.IsRecordNotFoundError(err) {
			err = errors.New("record not found")
		}

		return
	}

	rowAffected = db.RowsAffected

	return
}

func (c *Comment) DeletePostComments(db *gorm.DB, pid uint64) (rowAffected int64, err error) {

	db = db.Debug().Model(Comment{}).Where("post_id = ?", pid).Delete(Comment{})
	if db.Error != nil {
		if gorm.IsRecordNotFoundError(err) {
			err = errors.New("record not found")
		}

		return
	}

	rowAffected = db.RowsAffected

	return
}
