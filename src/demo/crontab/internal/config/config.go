package config

import (
	"github.com/IcecreamLee/goutils"
	"gopkg.in/ini.v1"
	"log"
)

var (
	Host               string
	Port               int
	User               string
	Password           string
	DBName             string
	CronTableName      string
	CronLogTableName   string
	ServiceName        string
	ServiceDisplayName string
	ServiceDescription string
)

func init() {
	iniConfig, err := ini.LoadSources(ini.LoadOptions{IgnoreInlineComment: true}, goutils.GetCurrentPath()+"crontab.ini")
	if err != nil {
		log.Fatalf("Fail to read file: %v", err)
	}
	section := iniConfig.Section("db")
	Host = section.Key("host").String()
	Port, _ = section.Key("port").Int()
	User = section.Key("user").String()
	Password = section.Key("password").String()
	DBName = section.Key("db_name").String()
	CronTableName = section.Key("cron_table_name").String()
	CronLogTableName = section.Key("cron_log_table_name").String()
	ServiceName = section.Key("service_name").String()
	ServiceDisplayName = section.Key("service_display_name").String()
	ServiceDescription = section.Key("service_description").String()
}
