package main

import (
	"fmt"
	"idworker/lib"
	"net/http"
)

func responseID(w http.ResponseWriter, r *http.Request) {
	id := lib.GetIDInstance().NextID()
	fmt.Fprintf(w, "%d", id) // 这个写入到w的是输出到客户端的
}

func main() {

	lib.Info("Hello")

	//http.HandleFunc("/", responseID)         // 设置访问的路由
	//err := http.ListenAndServe(":9999", nil) // 设置监听的端口
	//if err != nil {
	//	log.Fatal("ListenAndServe: ", err)
	//}
}
