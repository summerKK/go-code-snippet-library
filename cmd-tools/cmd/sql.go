package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/summerKK/go-code-snippet-library/cmd-tools/internal/sql2struct"
)

var username string
var password string

var host string
var charset string
var dbType string
var dbName string
var tableName string

var sqlCmd = &cobra.Command{
	Use:   "sql",
	Short: "sql转换和处理",
	Long:  "sql转换和处理",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

var structCmd = &cobra.Command{
	Use:   "struct",
	Short: "sql转换",
	Long:  "sql转换",
	Run: func(cmd *cobra.Command, args []string) {
		dbInfo := &sql2struct.DbInfo{
			DbType:   dbType,
			Host:     host,
			UserName: username,
			Password: password,
			Charset:  charset,
		}

		dbModel := sql2struct.NewDbModel(dbInfo)
		err := dbModel.Connect()
		if err != nil {
			log.Fatalf("dbModel.Connect error:%v", err)
		}
		columns, err := dbModel.GetColumns(dbName, tableName)
		if err != nil {
			log.Fatalf("dbModel.GetColumns error:%v", err)
		}

		template := sql2struct.NewStructTemplate()
		structColumns := template.AssemblyColumns(columns)
		err = template.Generate(tableName, structColumns)
		if err != nil {
			log.Fatalf("tempate.Generate error:%v", err)
		}
	},
}

func init() {
	sqlCmd.AddCommand(structCmd)

	structCmd.Flags().StringVarP(&username, "username", "u", "root", "数据库账号")
	structCmd.Flags().StringVarP(&password, "password", "p", "root", "数据库密码")
	structCmd.Flags().StringVarP(&host, "host", "H", "127.0.0.1", "数据库地址")
	structCmd.Flags().StringVarP(&charset, "charset", "c", "utf8", "数据库编码")
	structCmd.Flags().StringVarP(&dbType, "dbType", "D", "mysql", "数据库类型")
	structCmd.Flags().StringVarP(&dbName, "dbname", "d", "", "数据库名称")
	structCmd.Flags().StringVarP(&tableName, "tableName", "t", "", "数据库表名")
}
