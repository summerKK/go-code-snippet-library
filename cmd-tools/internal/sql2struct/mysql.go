package sql2struct

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type DbModel struct {
	DbEngine *sql.DB
	DbInfo   *DbInfo
}

type DbInfo struct {
	DbType   string
	Host     string
	UserName string
	Password string
	Charset  string
}

type TableColumn struct {
	ColumnName    string
	DataType      string
	IsNullable    string
	ColumnKey     string
	ColumnType    string
	ColumnComment string
}

var DbTypeToStructType = map[string]string{
	"int":        "int32",
	"tinyint":    "int8",
	"smallint":   "int",
	"mediumint":  "int64",
	"bigint":     "int64",
	"bit":        "int",
	"bool":       "bool",
	"enum":       "string",
	"set":        "string",
	"varchar":    "string",
	"char":       "string",
	"tinytext":   "string",
	"mediumtext": "string",
	"text":       "string",
	"longtext":   "string",
	"blob":       "string",
	"tinyblob":   "string",
	"mediumblob": "string",
	"longblob":   "string",
	"date":       "time.Time",
	"datetime":   "time.Time",
	"timestamp":  "time.Time",
	"time":       "time.Time",
	"float":      "float64",
	"double":     "float64",
}

func NewDbModel(info *DbInfo) *DbModel {
	return &DbModel{DbInfo: info}
}

func (m *DbModel) Connect() error {
	var err error
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s)/information_schema?charset=%s&parseTime=True&loc=Local",
		m.DbInfo.UserName,
		m.DbInfo.Password,
		m.DbInfo.Host,
		m.DbInfo.Charset,
	)
	m.DbEngine, err = sql.Open(m.DbInfo.DbType, dsn)

	return err
}

func (m *DbModel) GetColumns(dbName, tableName string) ([]*TableColumn, error) {
	sql := "select " + "COLUMN_NAME,DATA_TYPE,COLUMN_KEY,IS_NULLABLE,COLUMN_TYPE,COLUMN_COMMENT" +
		" from COLUMNS where TABLE_SCHEMA = ? and TABLE_NAME = ? "
	rows, err := m.DbEngine.Query(sql, dbName, tableName)
	if err != nil {
		return nil, err
	}
	if rows == nil {
		return nil, errors.New("没有数据")
	}
	defer rows.Close()

	var columns []*TableColumn
	for rows.Next() {
		var column TableColumn
		err := rows.Scan(&column.ColumnName, &column.DataType, &column.ColumnKey, &column.IsNullable, &column.ColumnType, &column.ColumnComment)
		if err != nil {
			return nil, err
		}

		columns = append(columns, &column)
	}

	return columns, nil
}
