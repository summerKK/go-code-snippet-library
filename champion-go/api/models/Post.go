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
	Content   string    `gorm:"size:255;not null;" json:"content"`
	Author    User      `json:"author"`
	AuthorID  uint32    `gorm:"not null" json:"author_id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (p *Post) Prepare() {
	p.ID = 0
	p.Title = html.EscapeString(strings.TrimSpace(p.Title))
	p.Content = html.EscapeString(strings.TrimSpace(p.Content))
	p.Author = User{}
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}

func (p *Post) Validate() error {

	if p.Title == "" {
		return errors.New("Required Title")
	}
	if p.Content == "" {
		return errors.New("Required Content")
	}
	if p.AuthorID < 1 {
		return errors.New("Required Author")
	}

	return nil
}

func (p *Post) SavePost(db *gorm.DB) (post *Post, err error) {
	post = &Post{}
	err = db.Debug().Model(Post{}).Create(&p).Error
	if err != nil {
		return
	}

	if p.ID > 0 {
		err = db.Debug().Model(User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
	}

	return
}

func (p *Post) FindAllPosts(db *gorm.DB) (posts []*Post, err error) {
	err = db.Debug().Model(Post{}).Limit(100).Find(&posts).Error
	if err != nil {
		return
	}
	if len(posts) > 0 {
		for _, post := range posts {
			err := db.Debug().Model(User{}).Where("id = ?", post.AuthorID).Take(&post.Author).Error
			if err != nil {
				return
			}
		}
	}

	return
}

func (p *Post) FindPostById(db *gorm.DB, pid uint64) (post *Post, err error) {
	post = &Post{}
	err = db.Debug().Model(Post{}).Where("id = ?", pid).Take(post).Error
	if err != nil {
		return
	}

	if p.ID > 0 {
		err = db.Debug().Model(User{}).Where("id = ?", p.AuthorID).Take(&post.Author).Error
	}

	return
}

func (p *Post) UpdateAPost(db *gorm.DB) (post *Post, err error) {
	post = &Post{}
	err = db.Debug().Model(Post{}).Where("id = ?", p.ID).Updates(
		Post{
			Title:     p.Title,
			Content:   p.Content,
			UpdatedAt: time.Now(),
		},
	).Error

	if err != nil {
		return
	}

	if p.ID > 0 {
		err = db.Debug().Model(User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
	}

	return
}

func (p *Post) DeleteAPost(db *gorm.DB, pid uint64, uid uint32) (rowAffected int64, err error) {
	db = db.Debug().Model(Post{}).Where("id = ? and author_id = ?", pid, uid).Take(&Post{}).Delete(&Post{})
	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			err = errors.New("post not found")
		}
		return
	}
	rowAffected = db.RowsAffected

	return
}
