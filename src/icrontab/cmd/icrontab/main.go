package main

import (
	"flag"
	"icrontab/internal/icrontab"
)

var ct *icrontab.ServiceManager

func main() {
	ct = icrontab.NewService()
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
