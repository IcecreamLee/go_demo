package main

import (
	"encoding/json"
	"fmt"
	"icecream/utils"
	"os"
)

var id *utils.IDGenServ

func main() {
	loadConfig()
	id.Run()
}

// 加载配置文件
func loadConfig() {
	id = utils.GetIDServInstance()
	file, _ := os.Open(utils.GetCurrentPath() + "conf.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	err := decoder.Decode(id)
	if err != nil {
		fmt.Println("Load Config Error:", err)
	}
}
