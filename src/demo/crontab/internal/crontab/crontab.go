package crontab

import (
	"Demo/crontab/internal/config"
	"Demo/crontab/internal/models"
	"fmt"
	"github.com/IcecreamLee/goutils"
	"github.com/robfig/cron/v3"
	"os/exec"
	"runtime"
	"time"
)

// 计划任务调度器
type Crontab struct {
	c        *cron.Cron
	cronJobs []*CronJob
	jobsMD5  string
}

var ct *Crontab

func New() *Crontab {
	if ct == nil {
		ct = new(Crontab)
		InitLogger()
	}
	return ct
}

func NewCrontab() *Crontab {
	return New()
}

// start 开启任务
func (c *Crontab) Start() {
	c.c = cron.New()
	var str string
	if c.cronJobs == nil {
		c.cronJobs = c.getCrons()
	}

	for _, curcron := range c.cronJobs {
		logger.Printf("add cron: %+v\n", curcron)
		entryID, _ := c.c.AddJob(curcron.Exp, curcron)
		curcron.NextExec = c.c.Entry(entryID).Schedule.Next(time.Now())
		curcron.CronId = int(entryID)
		curcron.Update()
		if c.jobsMD5 == "" {
			str = str + fmt.Sprintf("%d%s%s%s", curcron.Cron.ID, curcron.Cron.Exp, curcron.Cron.ExecType, curcron.Cron.ExecTarget)
		}
	}
	if str != "" {
		c.jobsMD5 = goutils.MD5Str([]byte(str))
	}
	go c.c.Start()

	// 每分钟获取一次数据，检查是都有更新，有更新则重新启动计划任务
	// 每小时打印一次goroutine数量，防止内存泄露
	mt := time.NewTicker(time.Minute)
	ht := time.NewTicker(time.Hour)
	for {
		select {
		case <-mt.C:
			if c.isChanged() {
				logger.Println("crontab is changed, restart...")
				c.Stop(false)
				go c.Start()
			}
		case <-ht.C:
			logger.Println("ticker, current goroutine num:", runtime.NumGoroutine())
		}
	}
}

// stop 停止任务
func (c *Crontab) Stop(isWaitDone bool) {
	ctx := c.c.Stop()
	for _, entry := range c.c.Entries() {
		c.c.Remove(entry.ID)
	}
	if isWaitDone {
		<-ctx.Done()
	}
}

// isCronUpdate 返回cron数据是否有更新
func (c *Crontab) isChanged() bool {
	newCrons := c.getCrons()
	var str string
	for _, curcron := range newCrons {
		str = str + fmt.Sprintf("%d%s%s%s", curcron.Cron.ID, curcron.Cron.Exp, curcron.Cron.ExecType, curcron.Cron.ExecTarget)
	}
	newCronsMD5 := goutils.MD5Str([]byte(str))
	if c.jobsMD5 == newCronsMD5 {
		return false
	}
	c.jobsMD5 = newCronsMD5
	c.cronJobs = newCrons
	return true
}

func (c Crontab) getCrons() []*CronJob {
	var jobs []*CronJob
	err := models.GetDB().Select(&jobs, `select id,cron_id,exp,exec_type,exec_target,last_exec,next_exec from `+config.CronTableName+` where is_delete = 0`)
	if err != nil {
		fmt.Printf("getCrons error: %s\n", err.Error())
	}
	return jobs
}

type CronJob struct {
	models.Cron
	cmd *exec.Cmd
}

// 计划任务运行
func (j *CronJob) Run() {
	entry := ct.c.Entry(cron.EntryID(j.Cron.CronId))
	if entry.Prev.Unix() > 0 {
		j.Cron.LastExec = entry.Prev
	}
	j.Cron.NextExec = entry.Next
	j.Cron.Update()

	if j.Cron.ExecTarget == "" {
		return
	}

	cronLog := &models.CronLog{
		CID:           j.Cron.ID,
		ExecStartTime: time.Now(),
		ExecType:      j.Cron.ExecType,
		ExecTarget:    j.Cron.ExecTarget,
	}

	if j.Cron.ExecType == "shell" {
		j.cmd = exec.Command(j.Cron.ExecTarget)
		err := j.cmd.Start()
		cronLog.PID = j.cmd.Process.Pid
		cronLog.Insert()

		err = j.cmd.Wait()
		if err == nil {
			cronLog.ExecStatus = 1
		} else {
			logger.Printf("[error] j run failed:%s\n", err)
			cronLog.ExecResult = "failed:" + err.Error()
		}
		cronLog.ExecEndTime = time.Now()
		cronLog.Update()
	}
}

// 任务强制停止
func (j *CronJob) Stop() {
	_ = j.cmd.Process.Kill()
}
