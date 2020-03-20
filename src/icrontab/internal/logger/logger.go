package logger

import (
	"github.com/IcecreamLee/goutils"
	"log"
	"os"
	"runtime/debug"
)

var (
	infoLog *log.Logger
	//debugLog *log.Logger
	errorLog *log.Logger
)

var (
	Info  func(v ...interface{})
	Infof func(format string, v ...interface{})
	//Debug  func(v ...interface{})
	//Debugf func(format string, v ...interface{})
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
	//debugLog = log.New(loggerFile, "[debug] ", log.LstdFlags)
	errorLog = log.New(loggerFile, "[error]", log.LstdFlags)
	Info = infoLog.Println
	Infof = infoLog.Printf
	//Debug = debugLog.Println
	//Debugf = debugLog.Printf
	Error = errorLog.Println
	Errorf = errorLog.Printf
}

// IfError 判断有错误则打印错误内容
func IfError(format string, err error, v ...interface{}) {
	if err != nil {
		v = append([]interface{}{err.Error()}, v...)
		Errorf(format, v...)
	}
}

// Recover 判断有崩溃则恢复奔溃并打印奔溃日志和调用栈信息
func Recover(msg string, handle func()) {
	if r := recover(); r != nil {
		Errorf("Recover: %s, ", r, msg)
		if handle != nil {
			handle()
		}
		Errorf(string(debug.Stack()))
	}
}
