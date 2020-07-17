package model

import (
	"github.com/jinzhu/gorm"
	"github.com/summerKK/go-code-snippet-library/blog-service/pkg/app"
)

type ArticleSwagger struct {
	List  []*Article
	Pager *app.Pager
}

type Article struct {
	*Model
	Title          string `json:"title"`
	Desc           string `json:"desc"`
	Content        string `json:"content"`
	ConverImageUrl string `json:"conver_image_url"`
	State          uint8  `json:"state"`
}

func (a Article) TableName() string {
	return "blog_article"
}

func (a Article) Count(db *gorm.DB) (int, error) {
	var count int
	if a.Title != "" {
		db = db.Where("title = ?", a.Title)
	}
	db = db.Where("state = ?", a.State)
	if err := db.Model(&a).Where("is_del = ?", 0).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (a Article) List(db *gorm.DB, pageSize, pageOffset int) ([]*Article, error) {
	var articles []*Article
	if pageSize >= 0 && pageOffset >= 0 {
		db = db.Offset(pageOffset).Limit(pageSize)
	}
	if a.Title != "" {
		db = db.Where("title = ?", a.Title)
	}
	db = db.Where("state = ?", a.State)

	if err := db.Where("is_del = ?", 0).Find(&articles).Error; err != nil {
		return nil, err
	}

	return articles, nil
}

func (a Article) Create(db *gorm.DB) error {
	return db.Create(&a).Error
}

func (a Article) Update(db *gorm.DB, vales interface{}) error {
	return db.Model(&a).Where("id = ? and is_del = ?", a.ID, 0).Updates(vales).Error
}

func (a Article) Delete(db *gorm.DB) error {
	return db.Where("id = ? and is_del = ?", a.ID, 0).Delete(&a).Error
}
