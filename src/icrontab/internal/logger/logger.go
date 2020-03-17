package logger

import (
	"github.com/IcecreamLee/goutils"
	"log"
	"os"
)

var (
	infoLog  *log.Logger
	errorLog *log.Logger
)

var (
	Info   func(v ...interface{})
	Infof  func(format string, v ...interface{})
	Error  func(v ...interface{})
	Errorf func(format string, v ...interface{})
)

func init() {
	var err error
	logFile := goutils.GetCurrentPath() + "icrontab.log"
	loggerFile, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	infoLog = log.New(loggerFile, "[info] ", log.LstdFlags)
	errorLog = log.New(loggerFile, "[error]", log.LstdFlags)
	Info = infoLog.Println
	Infof = infoLog.Printf
	Error = errorLog.Println
	Errorf = errorLog.Printf
}
