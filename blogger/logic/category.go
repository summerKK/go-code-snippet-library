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

func GetCategoryById(id int64) (category *model.Category, err error) {
	category = &model.Category{}
	category, err = db.CategoryById(id)
	return
}
