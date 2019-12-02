module demo/table_executor

go 1.13

require (
	github.com/go-sql-driver/mysql v1.4.1
	github.com/jmoiron/sqlx v1.2.0
	google.golang.org/appengine v1.6.5 // indirect
	gopkg.in/ini.v1 v1.51.0
	icecream v0.0.0-00010101000000-000000000000
)

replace icecream => ../../icecream
