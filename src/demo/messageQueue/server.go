package main

import (
	"fmt"
	"log"
	"net/http"
)

// startServer 启动一个web服务去接收消息队列发布请求
func startServer() {
	fmt.Println("start server")
	http.HandleFunc("/publish", publish)
	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		fmt.Println(err.Error())
	}
}

// publish 消息队列发布接收接口
func publish(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Printf("request params error: %s", err.Error())
		return
	}
	job := Message{
		Topic: r.Form.Get("topic"),
		Data:  r.Form.Get("data"),
	}
	MQSingleton().publish(job)
	_, _ = w.Write([]byte("success"))
}
