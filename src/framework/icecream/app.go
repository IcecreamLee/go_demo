package icecream

import (
	"log"
	"net/http"
)

func init() {
	NewApp()
}

var App *Icecream

func NewApp() *Icecream {
	App = new(Icecream)
	//App.IConfig.env = "dev"
	return App
}

func Run() *Icecream {
	http.HandleFunc("/", NewRequest)         //设置访问的路由
	err := http.ListenAndServe(":9090", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	return App
}

type Icecream struct {
	IConfig IConfig
}

type IConfig struct {
	Env     string
	Version string
}

//Appinit = > 加载配置，注册路由，启动服务
