package lib

import (
	"log"
	"os"
)

func Info(msg string)  {
	fileName := "/my_golang.log"
	logFile,err  := os.OpenFile(fileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	defer logFile.Close()
	if err != nil {
		log.Fatalln("open file error !")
	}
	// 创建一个日志对象
	debugLog := log.New(logFile, "[Info]", log.LstdFlags)
	debugLog.Println(msg)
}