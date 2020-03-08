package crontab

import (
	"github.com/IcecreamLee/goutils"
	"gopkg.in/ini.v1"
	"log"
)

var Conf *Config

type Config struct {
	Host             string
	Port             int
	User             string
	Password         string
	DBName           string
	CronTableName    string
	CronLogTableName string
}

func init() {
	Conf = &Config{}
	iniConfig, err := ini.LoadSources(ini.LoadOptions{IgnoreInlineComment: true}, goutils.GetCurrentPath()+"crontab.ini")
	if err != nil {
		log.Fatalf("Fail to read file: %v", err)
	}
	section := iniConfig.Section("db")
	Conf.Host = section.Key("Host").String()
	Conf.Port, _ = section.Key("Port").Int()
	Conf.User = section.Key("User").String()
	Conf.Password = section.Key("Password").String()
	Conf.DBName = section.Key("db_name").String()
	Conf.CronTableName = section.Key("cron_table_name").String()
	Conf.CronLogTableName = section.Key("cron_log_table_name").String()
}
