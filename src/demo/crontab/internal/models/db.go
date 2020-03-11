package models

import (
	"Demo/crontab/internal/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"
	"strconv"
)

var DB *sqlx.DB

func initDB() {
	var err error
	// user:password@tcp(localhost:3306)/dbname?params
	DB, err = sqlx.Connect("mysql", config.User+":"+config.Password+"@tcp("+config.Host+":"+strconv.Itoa(config.Port)+")/"+config.DBName+"?charset=utf8&parseTime=true&loc=Local")
	if err != nil {
		log.Fatalln(err)
	}
}

func GetDB() *sqlx.DB {
	if DB == nil {
		initDB()
	}
	return DB
}
