package icrontab

import (
	"github.com/kardianos/service"
	"icrontab/internal/config"
	"icrontab/internal/logger"
	"log"
	"os"
)

var sm *ServiceManager

type ServiceManager struct {
	ServiceName        string
	ServiceDisplayName string
	ServiceDescription string
	Service            service.Service
	icrontab           *ICrontab
}

// 接受命令行参数，已接受到返回true,否则返回false
func (s *ServiceManager) Manage(action string) {
	err := service.Control(sm.Service, action)
	if err != nil {
		logger.Infof("Valid actions: %q\n", service.ControlAction)
		logger.Error(err)
		os.Exit(1)
	}
}

// 运行服务
func (s *ServiceManager) Run() {
	err := sm.Service.Run()
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}
}

func NewService() *ServiceManager {
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
	icrontab *ICrontab
	exit     chan struct{}
}

// 服务运行
func (p *program) Start(s service.Service) error {
	if service.Interactive() {
		logger.Info("Running in terminal.")
	} else {
		logger.Info("Running under service manager.")
	}
	p.exit = make(chan struct{})

	// Start should not block. Do the actual work async.
	logger.Infof("I'm running %v.", service.Platform())

	p.icrontab = NewICrontab()
	go p.icrontab.Start()
	return nil
}

// 服务停止
func (p *program) Stop(s service.Service) error {
	// Any work in Stop should be quick, usually a few seconds at most.
	logger.Info("I'm Stopping...")
	close(p.exit)
	p.icrontab.Stop()
	logger.Info("I'm Stopped!")
	return nil
}
