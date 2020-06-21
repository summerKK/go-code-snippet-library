package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Post struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Title     string    `gorm:"size:255;not null;unique" json:"title"`
	Content   string    `gorm:"text;not null;" json:"content"`
	Author    *User     `json:"author"`
	AuthorID  uint32    `gorm:"not null" json:"author_id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (p *Post) Prepare() {

	p.Title = html.EscapeString(strings.TrimSpace(p.Title))
	p.Content = html.EscapeString(strings.TrimSpace(p.Content))
	p.Author = &User{}
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}

func (p *Post) Validate() map[string]string {

	var err error

	var errMsgList = make(map[string]string)

	if p.Title == "" {
		err = errors.New("Required Title")
		errMsgList["Required_title"] = err.Error()

	}
	if p.Content == "" {
		err = errors.New("Required Content")
		errMsgList["Required_content"] = err.Error()

	}
	if p.AuthorID < 1 {
		err = errors.New("Required Author")
		errMsgList["Required_author"] = err.Error()
	}

	return errMsgList
}

func (p *Post) SavePost(db *gorm.DB) (post *Post, err error) {

	p.Author = &User{}
	err = db.Debug().Model(Post{}).Create(&p).Error
	if err != nil {
		return
	}

	if p.ID > 0 {
		err = db.Debug().Model(User{}).Where("id = ?", p.AuthorID).Take(p.Author).Error
	}

	post = p

	return
}

func (p *Post) FindAllPosts(db *gorm.DB) (posts []*Post, err error) {

	err = db.Debug().Model(Post{}).Limit(100).Find(&posts).Error
	if err != nil {
		return
	}
	if len(posts) > 0 {
		for _, post := range posts {
			post.Author = &User{}
			err = db.Debug().Model(User{}).Where("id = ?", post.AuthorID).Take(post.Author).Error
			if err != nil {
				return
			}
		}
	}

	return
}

func (p *Post) FindPostById(db *gorm.DB, pid uint64) (post *Post, err error) {

	post = &Post{
		Author: &User{},
	}
	err = db.Debug().Model(Post{}).Where("id = ?", pid).Take(post).Error
	if err != nil {
		return
	}

	if p.ID > 0 {
		err = db.Debug().Model(User{}).Where("id = ?", p.AuthorID).Take(post.Author).Error
	}

	return
}

func (p *Post) UpdateAPost(db *gorm.DB) (post *Post, err error) {

	err = db.Debug().
		Model(Post{}).
		Where("id = ?", p.ID).
		Updates(
			Post{
				Title:     p.Title,
				Content:   p.Content,
				UpdatedAt: time.Now(),
			},
		).
		Error

	if err != nil {
		return
	}

	p.Author = &User{}
	if p.ID > 0 {
		err = db.Debug().Model(User{}).Where("id = ?", p.AuthorID).Take(p.Author).Error
	}

	post = p

	return
}

func (p *Post) DeleteAPost(db *gorm.DB, pid uint64, uid uint32) (rowAffected int64, err error) {

	db = db.Debug().Model(Post{}).Where("id = ? and author_id = ?", pid, uid).Delete(Post{})
	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			err = errors.New("post not found")
		}
		return
	}
	rowAffected = db.RowsAffected

	return
}

func (p *Post) FindUserPosts(db *gorm.DB, uid uint32) (posts []*Post, err error) {

	err = db.Debug().Model(Post{}).Where("author_id = ?", uid).Find(&posts).Error
	if err != nil {
		return
	}
	var user = &User{}
	err = db.Debug().Model(User{}).Where("id = ?", uid).Take(user).Error
	if err != nil {
		return
	}

	if len(posts) > 0 {
		for _, post := range posts {
			post.Author = user
		}
	}

	return
}

func (p *Post) DelteUserPosts(db *gorm.DB, uid uint32) (rowAffectd int64, err error) {

	db = db.Debug().Model(Post{}).Where("id = ?", uid).Delete(Post{})
	if db.Error != nil {
		err = db.Error
		return
	}
	rowAffectd = db.RowsAffected

	return
}
