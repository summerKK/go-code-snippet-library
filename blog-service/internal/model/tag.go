package model

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/summerKK/go-code-snippet-library/blog-service/pkg/app"
)

type TagSwagger struct {
	List  []*Tag
	Pager *app.Pager
}

type Tag struct {
	*Model
	Name  string `json:"name"`
	State uint8  `json:"state"`
}

func (t Tag) TableName() string {
	return "blog_tag"
}

func (t Tag) Count(db *gorm.DB) (int, error) {
	var count int
	if t.Name != "" {
		db = db.Where("name like ?", fmt.Sprintf("%%%s%%", t.Name))
	}
	db = db.Where("state = ?", t.State)
	if err := db.Model(&t).Where("is_del = ?", 0).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (t Tag) List(db *gorm.DB, pageOffset, pageSize int) ([]*Tag, error) {
	var tags []*Tag
	var err error
	if pageOffset >= 0 && pageSize > 0 {
		db = db.Offset(pageOffset).Limit(pageSize)
	}
	if t.Name != "" {
		db = db.Where("name like ?", fmt.Sprintf("%%%s%%", t.Name))
	}

	db = db.Where("state = ?", t.State)
	if err = db.Where("is_del = ?", 0).Find(&tags).Error; err != nil {
		return nil, err
	}

	return tags, nil
}

func (t Tag) Create(db *gorm.DB) error {
	return db.Create(&t).Error
}

func (t Tag) Update(db *gorm.DB, values interface{}) error {
	return db.Model(&t).Where("id = ? and is_del = ?", t.ID, 0).Updates(values).Error
}

func (t Tag) Delete(db *gorm.DB) error {
	return db.Where("id = ? and is_del = ?", t.ID, 0).Delete(&t).Error
}

func (t Tag) GetTag(db *gorm.DB) (*Tag, error) {
	var tag *Tag
	err := db.Model(&t).Where("id = ? and and state = ? and is_del = ?", t.ID, t.State, 0).First(&tag).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return tag, err
}
