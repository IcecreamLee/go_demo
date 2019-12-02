package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/ini.v1"
	"icecream/utils"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

var queryTableSql string
var execSqlForPerTable string

func main() {
	loadConfig()
	initDB()

	tables, err := utils.DBQuery(db, queryTableSql)
	fmt.Println(queryTableSql)
	if err != nil {
		log.Fatalf("Fail to query tables: %v", err)
		os.Exit(1)
	}

	fmt.Println("Table List:")
	var wg sync.WaitGroup
	tableNum := len(tables)
	convertedNum := 0
	for _, table := range tables {
		tableName := table["table"].(string)
		fmt.Println("    " + tableName)
		log2File("    " + tableName)
	}

	fmt.Println("\nExecuting...")
	log2File("Executing...")
	for _, table := range tables {
		wg.Add(1)
		t := table
		go func() {
			execSql := fmt.Sprintf(execSqlForPerTable, t["table"])
			_, err := db.Exec(execSql)
			if err != nil {
				log2File("    Table[%s] execute failed (%s): %v \n", t["table"], execSql, err)
			} else {
				log2File("    Table[%s] execute successfully (%s) \n", t["table"], execSql)
			}

			convertedNum++
			p := convertedNum * 50 / tableNum
			equalStr := strings.Repeat("=", p)
			space := strings.Repeat(" ", 50-p)
			fmt.Printf("\r [" + equalStr + ">" + space + "](" + fmt.Sprintf("%.1f", float64(convertedNum)/float64(tableNum)*100) + "/100%%)")
			wg.Done()
		}()
	}

	wg.Wait()
	log2File("All table is executed")
	fmt.Println("\n\nAll table is executed")
}

// sqlx
var db *sql.DB
var dbOnce sync.Once

func initDB() *sql.DB {
	var err error
	dbOnce.Do(func() {
		dataSource := dbConf.user + ":" + dbConf.password + "@tcp(" + dbConf.host + ":" + strconv.Itoa(dbConf.port) + ")/" + dbConf.dbName + "?charset=utf8&parseTime=True&loc=Local"
		db, err = sql.Open("mysql", dataSource)
		if err != nil {
			log.Fatalf("Fail to connect mysql: %v", err)
		}
	})
	return db
}

var dbConf dbConfig

type dbConfig struct {
	host     string
	port     int
	user     string
	password string
	dbName   string
}

func loadConfig() {
	config, err := ini.LoadSources(ini.LoadOptions{IgnoreInlineComment: true}, utils.GetCurrentPath()+"my.ini")
	if err != nil {
		log.Fatalf("Fail to read file: %v", err)
	}
	section := config.Section("db")
	dbConf.host = section.Key("host").String()
	dbConf.port, _ = section.Key("port").Int()
	dbConf.user = section.Key("user").String()
	dbConf.password = section.Key("password").String()
	dbConf.dbName = section.Key("db_name").String()
	section = config.Section("default")
	queryTableSql = section.Key("query_table_sql").String()
	execSqlForPerTable = section.Key("exec_sql_for_per_table").String()
	fmt.Printf("DB Config:\n    %+v \n\n", dbConf)
	log2File("DB Config: %+v", dbConf)
}

func log2File(format string, message ...interface{}) {
	utils.FileLogPrintf(utils.GetCurrentPath()+"run.log", format, message...)
}
