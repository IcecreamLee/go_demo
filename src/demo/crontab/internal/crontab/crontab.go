package crontab

import (
	"Demo/crontab/internal/config"
	"Demo/crontab/internal/models"
	"context"
	"fmt"
	"github.com/IcecreamLee/goutils"
	"github.com/robfig/cron/v3"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"time"
)

var (
	ct           *Crontab
	CMDProcesses map[int]*exec.Cmd
)

func init() {
	CMDProcesses = make(map[int]*exec.Cmd)
}

// 计划任务调度器
type Crontab struct {
	c          *cron.Cron
	cronJobs   []*CronJob
	jobsMD5    string
	httpServer *http.Server
}

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
	logger.Printf("start crontab...")
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

	go c.startHttpServer()

	// 每分钟获取一次数据，检查是都有更新，有更新则重新启动计划任务
	// 每小时打印一次goroutine数量，防止内存泄露
	mt := time.NewTicker(time.Minute)
	ht := time.NewTicker(time.Hour)
	for {
		select {
		case <-mt.C:
			if c.isChanged() {
				logger.Println("crontab is changed, restart...")
				c.restart()
			}
		case <-ht.C:
			logger.Println("ticker, current goroutine num:", runtime.NumGoroutine())
		}
	}
}

func (c *Crontab) startHttpServer() {
	if c.httpServer != nil {
		return
	}

	c.httpServer = &http.Server{Addr: ":" + config.ServicePort}

	// 停止子进程
	http.HandleFunc("/stop", func(writer http.ResponseWriter, request *http.Request) {
		_ = request.ParseForm()
		id := request.PostFormValue("id")
		cronLog := (&models.CronLog{ID: goutils.ToInt(id)}).Get("*")
		var err error
		cmd, ok := CMDProcesses[cronLog.PID]
		if ok {
			err = cmd.Process.Kill()
			delete(CMDProcesses, cronLog.PID)
			if err != nil {
				_, _ = fmt.Fprintln(writer, `{"code":1,"msg":"操作失败，`+err.Error()+`"}`)
				return
			}
		} else {
			_, _ = fmt.Fprintln(writer, `{"code":1,"msg":"操作失败，子进程不存在"}`)
			return
		}
		_, _ = fmt.Fprintln(writer, `{"code":0,"msg":"操作成功"}`)
	})

	// 重启任务
	http.HandleFunc("/restart", func(writer http.ResponseWriter, request *http.Request) {
		c.cronJobs = c.getCrons()
		c.restart()
		_, _ = fmt.Fprintln(writer, `{"code":0,"msg":"操作成功"}`)
	})

	http.HandleFunc("/ping", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = fmt.Fprintln(writer, "pong")
	})

	err := c.httpServer.ListenAndServe()
	if err != nil {
		log.Fatalln("start http service failed: ", err.Error())
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
		// 关闭http server
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		if err := c.httpServer.Shutdown(ctx); err != nil {
			logger.Fatal("Stop httpServer Failure: ", err)
		}
		logger.Println("Stopped httpServer")
	}
}

//
func (c *Crontab) restart() {
	c.Stop(false)
	go c.Start()
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
	err := models.GetDB().Select(&jobs, `select id,cron_id,exp,exec_type,exec_target,last_exec,next_exec from `+config.CronTableName+` where is_enable = 1 and is_delete = 0`)
	if err != nil {
		fmt.Printf("getCrons error: %s\n", err.Error())
	}
	return jobs
}

type CronJob struct {
	models.Cron
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
		cmd := exec.Command(j.Cron.ExecTarget)
		err := cmd.Start()
		cronLog.PID = cmd.Process.Pid
		cronLog.IsCrontab = 1
		cronLog.Insert()
		CMDProcesses[cronLog.PID] = cmd

		err = cmd.Wait()
		delete(CMDProcesses, cronLog.PID)
		if err != nil {
			cronLog.ExecResult = "failed:" + err.Error()
		}
		cronLog.ExecStatus = 1
		cronLog.ExecEndTime = time.Now()
		cronLog.Update()
	}
}
