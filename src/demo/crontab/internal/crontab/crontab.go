package crontab

import (
	"Demo/crontab/internal/models"
	"fmt"
	"github.com/IcecreamLee/goutils"
	"github.com/robfig/cron/v3"
	"os/exec"
	"runtime"
	"time"
)

// 计划任务调度器
type ICron struct {
	goCron   *cron.Cron
	crons    []*CronJob
	cronsMD5 string
}

var ic *ICron

func New() *ICron {
	if ic == nil {
		ic = &ICron{}
		InitLogger()
	}
	return ic
}

func NewCrontab() *ICron {
	return New()
}

// start 开启任务
func (ic *ICron) start() {
	ic.goCron = cron.New()
	var str string
	if ic.crons == nil {
		ic.crons = ic.getCrons()
	}
	for _, curcron := range ic.crons {
		logger.Printf("add cron: %+v\n", curcron)
		entryID, _ := ic.goCron.AddJob(curcron.Cron.Exp, curcron)
		curcron.Cron.NextExec = ic.goCron.Entry(entryID).Schedule.Next(time.Now())
		curcron.Cron.CronId = int(entryID)
		curcron.Cron.Update()
		if ic.cronsMD5 == "" {
			str = str + fmt.Sprintf("%d%s%s%s", curcron.Cron.ID, curcron.Cron.Exp, curcron.Cron.ExecType, curcron.Cron.ExecTarget)
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
func (ic *ICron) stop(isWaitDone bool) {
	ctx := ic.goCron.Stop()
	for _, entry := range ic.goCron.Entries() {
		ic.goCron.Remove(entry.ID)
	}
	if isWaitDone {
		<-ctx.Done()
	}
}

// isCronUpdate 返回cron数据是否有更新
func (ic *ICron) isChanged() bool {
	newCrons := ic.getCrons()
	var str string
	for _, curcron := range newCrons {
		str = str + fmt.Sprintf("%d%s%s", curcron.Cron.ID, curcron.Cron.Exp, curcron.Cron.ExecType, curcron.Cron.ExecTarget)
	}
	newCronsMD5 := goutils.MD5Str([]byte(str))
	if ic.cronsMD5 == newCronsMD5 {
		return false
	}
	ic.cronsMD5 = newCronsMD5
	ic.crons = newCrons
	return true
}

func (ic ICron) getCrons() []*CronJob {
	var crons []*CronJob
	err := models.GetDB().Select(&crons, `select * from `+Conf.CronTableName+` where is_delete = 0`)
	if err != nil {
		fmt.Printf("getCron error: %s\n", err.Error())
	}
	return crons
}

type CronJob struct {
	Cron *models.Cron
}

// 计划任务运行
func (cronJob *CronJob) Run() {
	entry := ic.goCron.Entry(cron.EntryID(cronJob.Cron.CronId))

	if entry.Prev.Unix() > 0 {
		cronJob.Cron.LastExec = entry.Prev
	}
	cronJob.Cron.NextExec = entry.Next
	cronJob.Cron.Update()

	if cronJob.Cron.ExecTarget == "" {
		return
	}

	cronLog := models.CronLog{
		CID:           cronJob.Cron.ID,
		ExecStartTime: time.Now(),
		ExecType:      cronJob.Cron.ExecType,
		ExecTarget:    cronJob.Cron.ExecTarget,
	}

	if cronJob.Cron.ExecType == "shell" {
		cmd := exec.Command(cronJob.Cron.ExecTarget)
		_, err := cmd.CombinedOutput()
		if err == nil {
			cronLog.ExecStatus = 1
		} else {
			logger.Printf("[error] cronJob run failed:%s\n", err)
			cronLog.ExecResult = "failed:" + err.Error()
		}
		cronLog.ExecEndTime = time.Now()
		cronLog.Insert()
	}
}
