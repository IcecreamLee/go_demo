package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/kardianos/service"
)

var logger service.Logger

type program struct{}

func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	Info("service start")
	go startIDService()
	return nil
}

func (p *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	Info("service stop")
	//stopIDService()
	<-time.After(time.Second * 13)
	return nil
}

func main() {
	svcConfig := &service.Config{
		Name:        "IceIDGeneratorService",
		DisplayName: "Ice ID Generator Service",
		Description: "This is an ID generator service.",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	if len(os.Args) > 1 {
		err = service.Control(s, os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	logger, err = s.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}
	err = s.Run()
	if err != nil {
		logger.Error(err)
	}
}

var idService *http.Server

func startIDService() {
	if idService == nil {
		idService = &http.Server{Addr: ":9999"}
		http.HandleFunc("/", idServiceResponse) // 设置访问的路由
	}

	err := idService.ListenAndServe()
	if err != nil {
		log.Fatal("StartIDService: ", err)
	}
}

func stopIDService() {
	//err := idService.Shutdown(context.Background())
	//if err != nil {
	//	log.Fatal("StopIDService: ", err)
	//}
}

func idServiceResponse(w http.ResponseWriter, r *http.Request) {
	id := GetIDInstance().NextID()
	fmt.Fprintf(w, "%d", id) // 这个写入到w的是输出到客户端的
}
