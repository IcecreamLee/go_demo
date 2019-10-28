package framework

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"sync"
)

func init() {
	_, _ = GetDB()
}

// gorm
var db *gorm.DB
var dbOnce sync.Once

func GetDB() (*gorm.DB, error) {
	var err error
	dbOnce.Do(func() {
		db, err = gorm.Open("mysql", "root:123456@tcp(localhost:3333)/test?charset=utf8&parseTime=True&loc=Local")
	})
	return db, err
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
