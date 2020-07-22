package model

import "github.com/jinzhu/gorm"

type Auth struct {
	*Model
	AppKey    string `json:"app_key"`
	AppSecret string `json:"app_secret"`
}

func (b Auth) TableName() string {
	return "blog_auth"
}

func (b Auth) Create(db *gorm.DB) error {
	return db.Create(&b).Error
}

func (b Auth) Get(db *gorm.DB) (*Auth, error) {
	auth := &Auth{}
	err := db.Where("app_key = ? and app_secret = ? and is_del = ?", b.AppKey, b.AppSecret, 0).First(&auth).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return auth, err
	}

	return auth, nil
}
