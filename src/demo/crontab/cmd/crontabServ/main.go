package main

import (
	"Demo/crontab/internal/crontab"
	"flag"
)

var ct *crontab.ServiceManager

func main() {
	ct = crontab.NewService()
	svcFlag := flag.String("s", "", "Control the system service.")
	flag.Parse()
	// 接受命令行参数标志控制服务
	action := *svcFlag
	if len(action) != 0 {
		ct.Manage(action)
		return
	}
	ct.Run()
}
