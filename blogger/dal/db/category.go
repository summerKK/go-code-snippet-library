package db

import (
	"github.com/jmoiron/sqlx"
	"summer/blogger/model"
)

func CategoryList(ids []int64) (list []*model.Category, err error) {
	sqlStr, args, err := sqlx.In("select * from category where id in (?)", ids)
	if err != nil {
		return nil, err
	}
	err = Db.Unsafe().Select(&list, sqlStr, args...)
	if err != nil {
		return
	}
	return
}

func AllCategoryList() (list []*model.Category, err error) {
	sqlStr := "select * from category order by category_no desc "
	err = Db.Unsafe().Select(&list, sqlStr)
	if err != nil {
		return
	}
	return
}

func CategoryById(id int64) (category *model.Category, err error) {
	sqlStr := "select * from category where id = ?"
	category = &model.Category{}
	err = Db.Unsafe().Get(category, sqlStr, id)
	return
}
