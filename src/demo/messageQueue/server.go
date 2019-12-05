package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

// startServer 启动一个web服务去接收消息队列发布请求
func startServer() {
	fmt.Println("Start server at localhost:" + strconv.Itoa(conf.Port) + "...")
	http.HandleFunc("/publish", publish)
	err := http.ListenAndServe(":"+strconv.Itoa(conf.Port), nil)
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

	priority, err := strconv.Atoi(r.Form.Get("priority"))
	if err != nil {
		priority = 0
	}

	executeTime, err := strconv.ParseInt(r.Form.Get("executeTime"), 10, 64)
	if err != nil {
		executeTime = time.Now().Unix()
	}

	job := Message{
		Topic:       r.Form.Get("topic"),
		Data:        r.Form.Get("data"),
		Priority:    priority,
		ExecuteTime: executeTime,
	}
	MQSingleton().publish(job)
	_, _ = w.Write([]byte("success"))
}
