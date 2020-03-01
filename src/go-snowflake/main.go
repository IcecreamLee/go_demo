package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/IcecreamLee/goutils"
	"github.com/kardianos/service"
	"log"
	"os"
)

func main() {
	svcFlag := flag.String("s", "", "Control the system service.")
	flag.Parse()

	// 服务定义
	svcConfig := &service.Config{
		Name:        "GoSnowflake",
		DisplayName: "Go Snowflake Service",
		Description: "This is an Go service that generate id.",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}

	// 接受命令行参数标志控制服务
	if len(*svcFlag) != 0 {
		err := service.Control(s, *svcFlag)
		if err != nil {
			log.Printf("Valid actions: %q\n", service.ControlAction)
			log.Fatal(err)
		}
		return
	}

	err = s.Run()
	if err != nil {
		goutils.LogInfo(goutils.GetCurrentPath()+"run.log", "GoSnowflake run failure: "+err.Error())
	}
}

// ID生成器http服务
var idServ *goutils.IDGenServ

// 加载配置文件
func init() {
	idServ = goutils.IDServSingleton()
	file, _ := os.Open(goutils.GetCurrentPath() + "conf.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	err := decoder.Decode(idServ)
	if err != nil {
		fmt.Println("Load Config Error:", err)
	}
}

// Program structures.
//  Define Start and Stop methods.
type program struct{}

// 服务运行
func (p *program) Start(s service.Service) error {
	goutils.LogInfo(goutils.GetCurrentPath()+"run.log", "GoSnowflake Starting...")
	// Start should not block. Do the actual work async.
	idServ = goutils.IDServSingleton()
	go idServ.Run()
	return nil
}

// 服务停止
func (p *program) Stop(s service.Service) error {
	// Any work in Stop should be quick, usually a few seconds at most.
	goutils.LogInfo(goutils.GetCurrentPath()+"run.log", "GoSnowflake Stopping...")
	idServ.Stop()
	goutils.LogInfo(goutils.GetCurrentPath()+"run.log", "GoSnowflake StStopped!")
	return nil
}
