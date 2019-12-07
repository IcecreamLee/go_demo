package main

import (
	"fmt"
	"github.com/IcecreamLee/goutils"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

// startServer 启动一个web服务去接收消息队列发布请求
func startServer() {
	fmt.Println("Start server at localhost:" + strconv.Itoa(conf.Port) + "...")
	http.HandleFunc("/publish", publish)
	http.HandleFunc("/testMsgCallback", testMsgCallback)
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
		ExecuteURL:  r.Form.Get("executeUrl"),
	}
	err = MQSingleton().publish(job)
	if err != nil {
		_, _ = w.Write([]byte("failure: " + err.Error()))
	}
	_, _ = w.Write([]byte("success"))
}

func testMsgCallback(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		println("Read the request body failure:", err.Error())
	}
	goutils.LogInfo(rootPath+"test.log", string(body))
	time.Sleep(time.Second)
	_, _ = w.Write(body)
}
