package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"
	"strconv"
)

var db *sqlx.DB

func initDB() {
	var err error
	// user:password@tcp(localhost:3306)/dbname?params
	db, err = sqlx.Connect("mysql", config.user+":"+config.password+"@tcp("+config.host+":"+strconv.Itoa(config.port)+")/"+config.dbName+"?charset=utf8&parseTime=true&loc=Local")
	if err != nil {
		log.Fatalln(err)
	}
}
