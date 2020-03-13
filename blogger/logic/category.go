package logic

import (
	"summer/blogger/dal/db"
	"summer/blogger/model"
)

func GetCategoryList() (list []*model.Category, err error) {
	list, err = db.AllCategoryList()
	if err != nil {
		return
	}
	return
}
