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
		return nil, err
	}
	return
}
