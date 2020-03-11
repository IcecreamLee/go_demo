package crontab

import (
	"Demo/crontab/internal/config"
	"github.com/kardianos/service"
	"log"
)

var sm *ServiceManager

type ServiceManager struct {
	ServiceName        string
	ServiceDisplayName string
	ServiceDescription string
	Service            service.Service
}

// 接受命令行参数，已接受到返回true,否则返回false
func (s *ServiceManager) Manage(action string) {
	err := service.Control(sm.Service, action)
	if err != nil {
		log.Printf("Valid actions: %q\n", service.ControlAction)
		log.Fatal(err)
	}
}

// 运行服务
func (s *ServiceManager) Run() {
	err := sm.Service.Run()
	if err != nil {
		log.Fatalln(err)
	}
}

func NewService() *ServiceManager {
	InitLogger()
	sm = &ServiceManager{
		ServiceName:        config.ServiceName,
		ServiceDisplayName: config.ServiceDisplayName,
		ServiceDescription: config.ServiceDescription,
	}

	// 服务定义
	svcConfig := &service.Config{
		Name:        sm.ServiceName,
		DisplayName: sm.ServiceDisplayName,
		Description: sm.ServiceDescription,
	}

	prg := &program{}
	var err error
	sm.Service, err = service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	return sm
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
	go NewCrontab().Start()
	return nil
}

// 服务停止
func (p *program) Stop(s service.Service) error {
	// Any work in Stop should be quick, usually a few seconds at most.
	logger.Println("I'm Stopping...")
	close(p.exit)
	ct.Stop(true)
	logger.Println("I'm Stopped!")
	return nil
}
