package main

import (
	"gopkg.in/ini.v1"
	"icecream/utils"
	"log"
)

var conf config

type config struct {
	Port               int
	MaximumConcurrency int
}

func init() {
	iniConfig, err := ini.LoadSources(ini.LoadOptions{IgnoreInlineComment: true}, utils.GetCurrentPath()+"my.ini")
	if err != nil {
		log.Fatalf("Fail to read file: %v", err)
	}
	section := iniConfig.Section("default")
	conf.Port, err = section.Key("port").Int()
	if err != nil {
		panic("MessageQueue service port is not config")
	}
	conf.MaximumConcurrency, err = section.Key("maximum_concurrency").Int()
	if err != nil {
		panic("MessageQueue MaximumConcurrency is not config")
	}
}
