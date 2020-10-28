package model

import "github.com/jinzhu/gorm"

type Album struct {
	*Model
	ArtistId uint32 `json:"artist_id"`
	Artist   Artist `json:"artist" gorm:"foreignKey:artist_id"`
	Name     string `json:"name"`
	Cover    string `json:"cover"`
}

func (a *Album) TableName() string {
	return "albums"
}

func (a *Album) Get(db *gorm.DB) (*Album, error) {
	var album Album
	err := db.Where("id = ?", a.ID).First(&album).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return &album, err
	}

	return &album, nil
}
