package models

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

type Like struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	UserID    uint32    `gorm:"not null" json:"user_id"`
	PostID    uint64    `gorm:"not null" json:"post_id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (l *Like) SaveLike(db *gorm.DB) (like *Like, err error) {

	like = &Like{}
	err = db.Debug().Model(Like{}).Where("user_id = ? and post_id = ?", l.UserID, l.PostID).Take(like).Error
	if err != nil {
		// 没有记录.给它添加记录
		if gorm.IsRecordNotFoundError(err) {
			err = db.Debug().Model(Like{}).Create(l).Take(like).Error
			if err != nil {
				return
			}
		}
	} else {
		// 记录以及存在
		err = errors.New("dobule like")
		return
	}

	return
}

func (l *Like) DeleteLike(db *gorm.DB) (deletedLike *Like, err error) {

	deletedLike = &Like{}
	db = db.Debug().Model(Like{}).Where("id = ?", l.ID).Take(deletedLike).Delete(deletedLike)
	if db.Error != nil {
		if gorm.IsRecordNotFoundError(err) {
			err = errors.New("record not found")
		}
		return
	}

	return
}

func (l *Like) GetLikesInfo(db *gorm.DB, pid uint64) (likeList []*Like, err error) {

	err = db.Debug().Model(Like{}).Where("post_id = ?", pid).Find(&likeList).Error

	return
}

func (l *Like) DeleteUserLikes(db *gorm.DB, uid uint32) (rowAffected int64, err error) {

	db = db.Debug().Model(Like{}).Where("user_id = ?", uid).Delete(&Like{})
	if db.Error != nil {
		if gorm.IsRecordNotFoundError(err) {
			err = errors.New("record not found")
		}
		return
	}

	rowAffected = db.RowsAffected

	return
}

func (l *Like) DeletePostLikes(db *gorm.DB, pid uint64) (rowAffected int64, err error) {

	db = db.Debug().Model(Like{}).Where("post_id = ?", pid).Delete(Like{})
	if db.Error != nil {
		if gorm.IsRecordNotFoundError(err) {
			err = errors.New("record not found")
		}
		return
	}

	rowAffected = db.RowsAffected

	return
}
