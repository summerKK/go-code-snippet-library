module github.com/summerKK/go-code-snippet-library/spark

require (
	github.com/gin-gonic/gin v1.5.0
	github.com/go-ozzo/ozzo-validation/v4 v4.1.0
	github.com/go-sql-driver/mysql v1.5.0
	github.com/jmoiron/sqlx v1.2.0
	github.com/sony/sonyflake v1.0.0
	github.com/summerKK/go-code-snippet-library/session v0.0.0-20200317154156-d88f1090bde7
	github.com/summerKK/go-code-snippet-library/trie v0.0.0-20200326134259-84a0e7c2f30c
)

replace (
	github.com/summerKK/go-code-snippet-library/session => ../session
	github.com/summerKK/go-code-snippet-library/trie => ../trie
)
