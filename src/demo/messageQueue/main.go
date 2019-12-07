package main

import (
	"github.com/IcecreamLee/goutils"
	"os"
	"os/signal"
)

var rootPath string

func main() {
	rootPath = goutils.GetCurrentPath()

	go handleExit()

	go MQSingleton().Run()

	startServer()
}

func handleExit() {
	c := make(chan os.Signal)
	// 监听所有信号
	signal.Notify(c)
	s := <-c
	goutils.FileLogPrintln(rootPath+"messageQueue", "exit", s)
	os.Exit(0)
}
