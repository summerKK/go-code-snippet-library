package model

import (
	"errors"

	"github.com/jinzhu/gorm"
	"github.com/summerKK/go-code-snippet-library/koel-api/pkg/util"
)

type Song struct {
	ID       string  `json:"id" gorm:"primary_key"`
	AlbumId  int     `json:"album_id"`
	Album    Album   `json:"album" gorm:"foreignKey:album_id"`
	ArtistId int     `json:"artist_id"`
	Artist   Artist  `json:"artist" gorm:"foreignKey:artist_id"`
	Title    string  `json:"title"`
	Length   float64 `json:"length"`
	Track    int     `json:"track"`
	Disc     int     `json:"disc"`
	Lyrics   string  `json:"lyrics"`
	Path     string  `json:"path"`
	Mtime    int     `json:"mtime"`
	TimeStruct
}

func (s *Song) TableName() string {
	return "songs"
}

//  钩子函数.插入前执行
func (s *Song) BeforeSave(tx *gorm.DB) error {
	if s.Path == "" {
		return errors.New("path can not be empty")
	}

	s.ID = util.EncodeMd5(s.Path)

	return nil
}

func (s *Song) List(db *gorm.DB, pageSize, pageOffset int) ([]*Song, error) {
	var songs []*Song
	if pageSize > 0 && pageOffset >= 0 {
		db = db.Offset(pageOffset).Limit(pageSize)
	}
	if s.Title != "" {
		db = db.Where("title = ?", s.Title)
	}

	if err := db.Find(&songs).Error; err != nil {
		return nil, err
	}

	return songs, nil
}

func (s *Song) All(db *gorm.DB) ([]*Song, error) {
	var songs []*Song
	if err := db.Find(&songs).Error; err != nil {
		return nil, err
	}

	return songs, nil
}
