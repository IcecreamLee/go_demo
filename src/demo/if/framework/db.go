package framework

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/jmoiron/sqlx"
	"icecream/utils"
	"sync"
)

func init() {
	connectDB()
}

// gorm
var GORM *gorm.DB
var DB *sql.DB
var dbOnce sync.Once

func connectDB() {
	var err error
	dbOnce.Do(func() {
		GORM, err = gorm.Open("mysql", "root:123456@tcp(localhost:3333)/test?charset=utf8&parseTime=True&loc=Local")
		if err != nil {
			panic(err)
		}
	})
	DB = GORM.DB()
}

// xorm
//var db *xorm.Engine
//var dbOnce sync.Once
//
//func GetDB() (*xorm.Engine, error) {
//	var err error
//	dbOnce.Do(func() {
//		db, err = xorm.NewEngine("mysql", "root:123456@tcp(localhost:3333)/test?charset=utf8&parseTime=True&loc=Local")
//	})
//	return db, err
//}

// sqlx
//var db *sqlx.DB
//var dbOnce sync.Once
//
//func GetDB() (*sqlx.DB, error) {
//	var err error
//	dbOnce.Do(func() {
//		db, err = sqlx.Open("mysql", "root:123456@tcp(localhost:3333)/test?charset=utf8&parseTime=True&loc=Local")
//	})
//	return db, err
//}

// 数据库查询多行数据，返回一个map键值对数组
// db: 传入db连接实例指针, query: 传入查询语句, args: 传入查询参数
func DBQuery(db *sql.DB, query string, args ...interface{}) ([]map[string]interface{}, error) {
	return utils.DBQuery(db, query, args...)
}

// 数据库查询单行数据，返回一个map键值对
// db: 传入db连接实例指针, query: 传入查询语句, args: 传入查询参数
func DBQueryRow(db *sql.DB, query string, args ...interface{}) (map[string]interface{}, error) {
	return utils.DBQueryRow(db, query, args...)
}

// 数据库查询多行数据，返回一个map键值对数组
// db: 传入sqlx.db连接实例指针, query: 传入查询语句, args: 传入查询参数
func XDBQuery(db *sqlx.DB, query string, args ...interface{}) ([]map[string]interface{}, error) {
	return utils.XDBQuery(db, query, args...)
}

// 数据库查询单行数据，返回一个map键值对
// db: 传入sqlx.db连接实例指针, query: 传入查询语句, args: 传入查询参数
func XDBQueryRow(db *sqlx.DB, query string, args ...interface{}) (map[string]interface{}, error) {
	return utils.XDBQueryRow(db, query, args...)
}

// MapScan将sql.Rows当前行数据赋值给dest
func MapScan(r *sql.Rows, dest map[string]interface{}) error {
	return utils.MapScan(r, dest)
}
