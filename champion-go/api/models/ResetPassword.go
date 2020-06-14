package models

import (
	"html"

	"strings"

	"github.com/jinzhu/gorm"
)

type ResetPassword struct {
	gorm.Model
	Email string `gorm:"size:100;not null;" json:"email"`
	Token string `gorm:"size:255;not null;" json:"token"`
}

func (r *ResetPassword) Prepare() {

	r.Token = html.EscapeString(strings.TrimSpace(r.Token))
	r.Email = html.EscapeString(strings.TrimSpace(r.Email))
}

func (r *ResetPassword) SaveDatails(db *gorm.DB) (*ResetPassword, error) {

	var err error
	err = db.Debug().Create(&r).Error
	if err != nil {
		return &ResetPassword{}, err
	}
	return r, nil
}

func (r *ResetPassword) DeleteDatails(db *gorm.DB) (int64, error) {

	db = db.Debug().Model(&ResetPassword{}).Where("id = ?", r.ID).Take(&ResetPassword{}).Delete(&ResetPassword{})

	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
