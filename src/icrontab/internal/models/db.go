package models

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"icrontab/internal/config"
	"icrontab/internal/logger"
	"os"
	"strconv"
)

var DB *sqlx.DB

func init() {
	GetDB()
}

func initDB() {
	var err error
	// user:password@tcp(localhost:3306)/dbname?params
	DB, err = sqlx.Connect("mysql", config.User+":"+config.Password+"@tcp("+config.Host+":"+strconv.Itoa(config.Port)+")/"+config.DBName+"?charset=utf8&parseTime=true&loc=Local")
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}
}

func GetDB() *sqlx.DB {
	if DB == nil {
		initDB()
	}
	return DB
}
