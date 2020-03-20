module github.com/summerKK/go-code-snippet-library/spark

require (
	github.com/DeanThompson/ginpprof v0.0.0-20190408063150-3be636683586
	github.com/gin-gonic/gin v1.5.0
	github.com/go-sql-driver/mysql v1.5.0
	github.com/jmoiron/sqlx v1.2.0
	github.com/sony/sonyflake v1.0.0
	github.com/summerKK/go-code-snippet-library/session v0.0.0-20200317154156-d88f1090bde7
)

replace github.com/summerKK/go-code-snippet-library/session => ../session
