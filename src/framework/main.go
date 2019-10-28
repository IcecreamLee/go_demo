package main

import (
	"fmt"
	"framework/icecream"
	_ "framework/routers"
	"log"
	"net/http"
)

func main() {
	app := icecream.Run()
	//str := Test()
	//fmt.Printf("%s", str)
	fmt.Printf("%#v", app)
	return

	http.HandleFunc("/", sayhelloName)       //设置访问的路由
	err := http.ListenAndServe(":9090", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func sayhelloName(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello astaxie!") //这个写入到w的是输出到客户端的
}
