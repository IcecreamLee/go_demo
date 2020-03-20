package models

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"icrontab/internal/config"
	"icrontab/internal/logger"
	"os"
	"strconv"
	"strings"
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

func UpdateWithMap(table string, data map[string]interface{}, condition map[string]interface{}) error {
	var updates []string
	var args []interface{}
	var conditions []string
	for key, value := range data {
		updates = append(updates, key+"=?")
		args = append(args, value)
	}
	for key, value := range condition {
		conditions = append(conditions, key+"=?")
		args = append(args, value)
	}

	sql := `UPDATE ` + table + ` SET ` + strings.Join(updates, ",") + ` WHERE ` + strings.Join(conditions, " AND ")
	_, err := GetDB().Exec(sql, args...)
	return err
}

func InsertWithMap(table string, data map[string]interface{}) error {
	var fields []string
	var placeholders []string
	var args []interface{}
	for key, value := range data {
		fields = append(fields, key)
		placeholders = append(placeholders, "?")
		args = append(args, value)
	}
	sql := "INSERT INTO " + table + " (" + strings.Join(fields, ",") + ") VALUES (" + strings.Join(placeholders, ",") + ")"
	_, err := GetDB().Exec(sql, args...)
	return err
}
