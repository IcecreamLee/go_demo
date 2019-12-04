package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"time"
)

func main() {

	go handleExit()

	go printGoroutineNum()

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

func printGoroutineNum() {
	for {
		fmt.Println("Goroutine num:", runtime.NumGoroutine())
		time.Sleep(time.Second * 3)
	}
}
