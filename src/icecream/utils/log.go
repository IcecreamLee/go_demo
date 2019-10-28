package utils

import (
	"log"
	"os"
)

func LogInfo(logFilePath string, msg string) {
	logFile, err := os.OpenFile(logFilePath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	defer logFile.Close()
	if err != nil {
		log.Fatalln("open file error !")
	}
	// 创建一个日志对象
	debugLog := log.New(logFile, "", log.LstdFlags)
	debugLog.Println(msg)
}
