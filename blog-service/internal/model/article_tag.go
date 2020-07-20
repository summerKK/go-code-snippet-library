package model

import "github.com/jinzhu/gorm"

type ArticleTag struct {
	*Model
	TagId     uint32 `json:"tag_id"`
	ArticleId uint32 `json:"article_id"`
}

func (a ArticleTag) TableName() string {
	return "blog_article_tag"
}

func (a ArticleTag) GetByAID(db *gorm.DB) (*ArticleTag, error) {
	var articleTag ArticleTag
	err := db.Where("article_id = ? and is_del = ?", a.ArticleId, 0).First(&articleTag).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &articleTag, nil
}

func (a ArticleTag) ListByTID(db *gorm.DB) ([]*ArticleTag, error) {
	var list []*ArticleTag
	err := db.Where("tag_id = ? and is_del = ?", a.TagId, 0).Find(&list).Error

	return list, err
}

func (a ArticleTag) ListByAIDs(db *gorm.DB, aIds []uint32) ([]*ArticleTag, error) {
	var list []*ArticleTag
	err := db.Where("article_id in (?) and is_del = ?", aIds, 0).Find(&list).Error

	return list, err
}

func (a ArticleTag) Create(db *gorm.DB) error {
	return db.Create(&a).Error
}

func (a ArticleTag) UpdateOne(db *gorm.DB, values interface{}) error {
	err := db.Model(&a).Where("article_id = ? and is_del = ?", a.ArticleId, 0).Limit(1).Updates(values).Error

	return err
}

func (a ArticleTag) Delete(db *gorm.DB) error {
	err := db.Where("id = ? and is_del = ?", a.ID, 0).Delete(&a).Error

	return err
}

func (a ArticleTag) DeleteOne(db *gorm.DB) error {
	err := db.Where("article_id = ? and is_del = ?", a.ArticleId, 0).Limit(1).Delete(&a).Error

	return err
}
