package main

import (
	"flag"
	"github.com/IcecreamLee/goutils"
	"github.com/kardianos/service"
	"gopkg.in/ini.v1"
	"log"
	"os"
)

var (
	c          *iCron
	logFile    string
	loggerFile *os.File
	logger     *log.Logger
	config     Config
)

func main() {
	defer loggerFile.Close()

	svcFlag := flag.String("s", "", "Control the system service.")
	flag.Parse()

	// 服务定义
	svcConfig := &service.Config{
		Name:        "GoCrontabService",
		DisplayName: "Go CronTab Service",
		Description: "This is an Go service that run cron jobs.",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}

	// 接受命令行参数标志控制服务
	if len(*svcFlag) != 0 {
		err := service.Control(s, *svcFlag)
		if err != nil {
			log.Printf("Valid actions: %q\n", service.ControlAction)
			log.Fatal(err)
		}
		return
	}

	err = s.Run()
	if err != nil {
		logger.Println(err)
	}
}

func init() {
	logFile = goutils.GetCurrentPath() + "go_crontab.log"
	initConfig()
	initDB()
	initLogger()
}

func initLogger() {
	var err error
	loggerFile, err = os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	logger = log.New(loggerFile, "", log.LstdFlags)
}

type Config struct {
	host             string
	port             int
	user             string
	password         string
	dbName           string
	cronTableName    string
	cronLogTableName string
}

func initConfig() {
	iniConfig, err := ini.LoadSources(ini.LoadOptions{IgnoreInlineComment: true}, goutils.GetCurrentPath()+"crontab.ini")
	if err != nil {
		log.Fatalf("Fail to read file: %v", err)
	}
	section := iniConfig.Section("db")
	config.host = section.Key("host").String()
	config.port, _ = section.Key("port").Int()
	config.user = section.Key("user").String()
	config.password = section.Key("password").String()
	config.dbName = section.Key("db_name").String()
	config.cronTableName = section.Key("cron_table_name").String()
	config.cronLogTableName = section.Key("cron_log_table_name").String()
}

// Program structures.
//  Define Start and Stop methods.
type program struct {
	exit chan struct{}
}

// 服务运行
func (p *program) Start(s service.Service) error {
	if service.Interactive() {
		logger.Println("Running in terminal.")
	} else {
		logger.Println("Running under service manager.")
	}
	p.exit = make(chan struct{})

	// Start should not block. Do the actual work async.
	logger.Printf("I'm running %v.", service.Platform())
	c = new(iCron)
	go c.start()
	return nil
}

// 服务停止
func (p *program) Stop(s service.Service) error {
	// Any work in Stop should be quick, usually a few seconds at most.
	logger.Println("I'm Stopping...")
	close(p.exit)
	c.stop(true)
	logger.Println("I'm Stopped!")
	return nil
}
