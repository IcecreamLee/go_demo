package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type Job struct {
	Topic        string
	Data         string
	Priority     int
	Status       int
	Result       int
	ExecuteTime  int64
	ExecutedTime int64
}

type JobStack struct {
	Jobs []Job
	Lock sync.RWMutex
}

type Daemon struct {
	JobStack           JobStack
	MaximumConcurrency int
}

func (d *Daemon) publish(job Job) {
	d.JobStack.Lock.Lock()
	defer d.JobStack.Lock.Unlock()
	d.JobStack.Jobs = append(d.JobStack.Jobs, job)
}

func (d *Daemon) consume() (Job, error) {
	d.JobStack.Lock.Lock()
	defer d.JobStack.Lock.Unlock()
	if len(d.JobStack.Jobs) == 0 {
		return Job{}, errors.New("no consume")
	}
	job := d.JobStack.Jobs[0]
	d.JobStack.Jobs = d.JobStack.Jobs[1:]
	return job, nil
}

func (d *Daemon) Run() {
	fmt.Println("Daemon running...")
	var wg sync.WaitGroup
	for {
		wg.Add(d.MaximumConcurrency)
		for i := 0; i < d.MaximumConcurrency; i++ {
			go func() {
				job, err := d.consume()
				if err == nil {
					fmt.Println("execute job:", job)
					time.Sleep(time.Second)
					fmt.Println("executed job:", job)
				}
				wg.Done()
			}()
		}
		wg.Wait()
	}
}

var daemonOnce sync.Once
var daemon *Daemon

func DaemonSingleton() *Daemon {
	daemonOnce.Do(func() {
		daemon = &Daemon{MaximumConcurrency: 3}
	})
	return daemon
}
