package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
)

var d *Daemon

func main() {

	go handleExit()

	d = DaemonSingleton()
	go d.Run()

	fmt.Println("start server")
	http.HandleFunc("/in", inQueue)
	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func inQueue(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Printf("request params error: %s", err.Error())
		return
	}
	job := Job{
		Topic: r.Form.Get("topic"),
		Data:  r.Form.Get("data"),
	}
	d.publish(job)
}

func handleExit() {
	c := make(chan os.Signal)
	//监听所有信号
	signal.Notify(c)
	s := <-c
	fmt.Println("exit", s)
	os.Exit(1)
}
