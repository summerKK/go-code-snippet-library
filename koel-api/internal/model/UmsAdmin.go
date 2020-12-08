package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

const (
	UmsAdminStatusValid int8 = iota
	UmsAdminStatusInvalidValid
)

type UmsAdmin struct {
	Id         uint      `json:"id" gorm:"primaryKey"`
	Username   string    `json:"username"`
	Password   string    `json:"-"`
	Icon       string    `json:"icon"`
	Email      string    `json:"email"`
	NickName   string    `json:"nick_name"`
	Note       string    `json:"note"`
	CreateTime time.Time `json:"create_time"`
	LoginTime  time.Time `json:"login_time"`
	Status     uint8     `json:"status"`
}

func (a *UmsAdmin) TableName() string {
	return "ums_admin"
}

func (a *UmsAdmin) GetUserByName(db *gorm.DB) (*UmsAdmin, error) {
	var user UmsAdmin
	err := db.Where("username = ?", a.Username).First(&user).Error

	return &user, err
}
