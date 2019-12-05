package main

import (
	"fmt"
	"os"
	"os/signal"
)

func main() {

	go handleExit()

	go MQSingleton().Run()

	startServer()
}

func handleExit() {
	c := make(chan os.Signal)
	// 监听所有信号
	signal.Notify(c)
	s := <-c
	fmt.Println("exit", s)
	os.Exit(0)
}
