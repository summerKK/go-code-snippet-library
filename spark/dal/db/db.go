package db

import "github.com/jmoiron/sqlx"
import _ "github.com/go-sql-driver/mysql"

var Db *sqlx.DB

func Init(dns string) error {
	var err error
	Db, err = sqlx.Open("mysql", dns)
	if err != nil {
		return err
	}

	err = Db.Ping()
	if err != nil {
		return err
	}

	Db.SetMaxOpenConns(100)
	Db.SetMaxIdleConns(16)

	return nil
}
