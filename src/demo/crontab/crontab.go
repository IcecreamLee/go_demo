package main

import (
	"fmt"
	"github.com/IcecreamLee/goutils"
	"github.com/robfig/cron/v3"
	"runtime"
	"time"
)

// 计划任务调度器
type iCron struct {
	goCron   *cron.Cron
	crons    []*Cron
	cronsMD5 string
}

// start 开启任务
func (ic *iCron) start() {
	ic.goCron = cron.New()
	var str string
	if ic.crons == nil {
		ic.crons = ic.getCrons()
	}
	for _, curcron := range ic.crons {
		logger.Printf("add cron: %+v\n", curcron)
		entryID, _ := ic.goCron.AddJob(curcron.Exp, curcron)
		curcron.NextExec = ic.goCron.Entry(entryID).Schedule.Next(time.Now())
		curcron.CronId = int(entryID)
		curcron.Update()
		if ic.cronsMD5 == "" {
			str = str + fmt.Sprintf("%d%s%s%s", curcron.ID, curcron.Exp, curcron.ExecType, curcron.ExecTarget)
		}
	}
	if str != "" {
		ic.cronsMD5 = goutils.MD5Str([]byte(str))
	}
	go ic.goCron.Start()

	// 每分钟获取一次数据，检查是都有更新，有更新则重新启动计划任务
	// 每小时打印一次goroutine数量，防止内存泄露
	mt := time.NewTicker(time.Minute)
	ht := time.NewTicker(time.Hour)
	for {
		select {
		case <-mt.C:
			if ic.isChanged() {
				logger.Println("crontab is changed, restart...")
				ic.stop(false)
				go ic.start()
			}
		case <-ht.C:
			logger.Println("ticker, current goroutine num:", runtime.NumGoroutine())
		}
	}
}

// stop 停止任务
func (ic *iCron) stop(isWaitDone bool) {
	ctx := ic.goCron.Stop()
	for _, entry := range ic.goCron.Entries() {
		ic.goCron.Remove(entry.ID)
	}
	if isWaitDone {
		<-ctx.Done()
	}
}

// isCronUpdate 返回cron数据是否有更新
func (ic *iCron) isChanged() bool {
	newCrons := ic.getCrons()
	var str string
	for _, curcron := range newCrons {
		str = str + fmt.Sprintf("%d%s%s", curcron.ID, curcron.Exp, curcron.ExecType, curcron.ExecTarget)
	}
	newCronsMD5 := goutils.MD5Str([]byte(str))
	if ic.cronsMD5 == newCronsMD5 {
		return false
	}
	ic.cronsMD5 = newCronsMD5
	ic.crons = newCrons
	return true
}

func (ic iCron) getCrons() []*Cron {
	var crons []*Cron
	err := db.Select(&crons, `select * from `+config.cronTableName+` where is_delete = 0`)
	if err != nil {
		fmt.Printf("getCron error: %s\n", err.Error())
	}
	return crons
}
