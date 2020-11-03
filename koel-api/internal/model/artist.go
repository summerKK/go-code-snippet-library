package model

import "github.com/jinzhu/gorm"

type Artist struct {
	Model
	Name  string `json:"name"`
	Image string `json:"image"`
}

func (a *Artist) TableName() string {
	return "artists"
}

func (a *Artist) Get(db *gorm.DB) (*Artist, error) {
	var artist Artist
	err := db.Where("id = ?", a.ID).First(&artist).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return &artist, err
	}

	return &artist, nil
}

func (a *Artist) First(db *gorm.DB) (*Artist, error) {
	var artist Artist
	err := db.First(&artist).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return &artist, err
	}

	return &artist, nil
}
