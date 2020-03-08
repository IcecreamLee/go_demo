package crontab

import (
	"flag"
	"github.com/kardianos/service"
	"log"
)

type Service struct {
	Name string
}

func (s *Service) Manage() {

}

func NewService() {
	svcFlag := flag.String("s", "", "Control the system service.")
	flag.Parse()

	// 服务定义
	svcConfig := &service.Config{
		Name:        "GoCrontabService",
		DisplayName: "Go CronTab Service",
		Description: "This is an Go service that run cron jobs.",
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
		logger.Println(err)
	}
}

// Program structures.
//  Define Start and Stop methods.
type program struct {
	exit chan struct{}
}

// 服务运行
func (p *program) Start(s service.Service) error {
	if service.Interactive() {
		logger.Println("Running in terminal.")
	} else {
		logger.Println("Running under service manager.")
	}
	p.exit = make(chan struct{})

	// Start should not block. Do the actual work async.
	logger.Printf("I'm running %v.", service.Platform())
	go NewCrontab().start()
	return nil
}

// 服务停止
func (p *program) Stop(s service.Service) error {
	// Any work in Stop should be quick, usually a few seconds at most.
	logger.Println("I'm Stopping...")
	close(p.exit)
	ic.stop(true)
	logger.Println("I'm Stopped!")
	return nil
}
