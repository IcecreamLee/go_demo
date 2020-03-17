package icrontab

import (
	"context"
	"fmt"
	"github.com/IcecreamLee/goutils"
	"github.com/robfig/cron/v3"
	"icrontab/internal/config"
	"icrontab/internal/logger"
	"icrontab/internal/models"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"
)

var ic *ICrontab

// ICrontab
type ICrontab struct {
	Scheduler  *Scheduler
	httpServer *http.Server
}

func New() *ICrontab {
	ic = &ICrontab{
		Scheduler: &Scheduler{
			c:              cron.New(),
			cronJobs:       nil,
			childProcesses: make(map[int]*exec.Cmd),
		},
		httpServer: &http.Server{Addr: ":" + config.ServicePort},
	}
	return ic
}

func NewICrontab() *ICrontab {
	return New()
}

// 启动ICrontab
func (i *ICrontab) Start() {
	logger.Info("start icrontab")

	// 启动任务调度器
	go i.Scheduler.Start()

	// 停止子进程
	http.HandleFunc("/stop", func(writer http.ResponseWriter, request *http.Request) {
		_ = request.ParseForm()
		id := request.PostFormValue("id")
		msg := i.Scheduler.StopChildProcess(goutils.ToInt(id))
		_, _ = fmt.Fprintln(writer, msg)
	})

	// 重启任务
	http.HandleFunc("/restart", func(writer http.ResponseWriter, request *http.Request) {
		i.Scheduler.Restart()
		_, _ = fmt.Fprintln(writer, `{"code":0,"msg":"操作成功"}`)
	})

	http.HandleFunc("/ping", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = fmt.Fprintln(writer, "pong")
	})

	err := i.httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		logger.Error("start http service failed: ", err.Error())
		os.Exit(1)
	}
}

// 停止ICrontab
func (i *ICrontab) Stop() {
	// 停止计划任务调度器
	ctx := i.Scheduler.Stop()
	<-ctx.Done()
	// 关闭http server
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := i.httpServer.Shutdown(ctx); err != nil {
		logger.Error("Stop httpServer Failure: ", err)
		os.Exit(1)
	}
	logger.Info("Stopped httpServer")
}

// 任务调度器
type Scheduler struct {
	c              *cron.Cron        // go crontab
	cronJobs       []*CronJob        // 任务配置列表
	crontabMD5     string            // 任务配置MD5
	childProcesses map[int]*exec.Cmd // 子进程
	minuteTicker   *time.Ticker
	hourTicker     *time.Ticker
}

// 启动任务调度器
func (s *Scheduler) Start() {
	var crontabs []*models.Crontab
	if s.cronJobs == nil {
		crontabs = models.GetCrons("id,exp,exec_type,exec_target,last_exec,next_exec")
	}
	str := ""
	for index, curCrontab := range crontabs {
		logger.Infof("add cron: %+v\n", curCrontab)
		s.cronJobs = append(s.cronJobs, &CronJob{s, curCrontab})
		entryID, _ := s.c.AddJob(curCrontab.Exp, s.cronJobs[index])
		curCrontab.NextExec = s.c.Entry(entryID).Schedule.Next(time.Now())
		curCrontab.CronId = int(entryID)
		curCrontab.Update()
		if s.crontabMD5 == "" {
			str = str + fmt.Sprintf("%d%s%s%s", curCrontab.ID, curCrontab.Exp, curCrontab.ExecType, curCrontab.ExecTarget)
		}
	}
	if str != "" {
		s.crontabMD5 = goutils.MD5Str([]byte(str))
	}
	go s.c.Start()

	go func() {
		if s.minuteTicker != nil {
			return
		}
		// 每分钟获取一次数据，检查是都有更新，有更新则重新启动计划任务
		// 每小时打印一次goroutine数量，防止内存泄露
		s.minuteTicker = time.NewTicker(time.Minute)
		s.hourTicker = time.NewTicker(time.Hour)
		for {
			select {
			case <-s.minuteTicker.C:
				if s.isChanged() {
					logger.Info("crontab is changed, restart...")
					s.Restart()
				}
			case <-s.hourTicker.C:
				logger.Info("ticker, current goroutine num:", runtime.NumGoroutine())
			}
		}
	}()
}

// 停止任务调度器
func (s *Scheduler) Stop() context.Context {
	ctx := s.c.Stop()
	for _, entry := range s.c.Entries() {
		s.c.Remove(entry.ID)
	}
	return ctx
}

// 停止子进程
func (s *Scheduler) StopChildProcess(id int) string {
	cronLog := (&models.CronLog{ID: id}).Get("*")
	var err error
	cmd, ok := s.childProcesses[cronLog.PID]
	if ok {
		err = cmd.Process.Kill()
		delete(s.childProcesses, cronLog.PID)
		if err != nil {
			return `{"code":1,"msg":"操作失败，` + err.Error() + `"}`
		}
	} else {
		return `{"code":1,"msg":"操作失败，子进程不存在"}`
	}
	return `{"code":0,"msg":"操作成功"}`
}

// 重启任务调度器
func (s *Scheduler) Restart() {
	logger.Info("restart scheduler...")
	s.Stop()
	s.cronJobs = nil
	go s.Start()
}

// isCronUpdate 返回cron数据是否有更新
func (s *Scheduler) isChanged() bool {
	newCronJobs := models.GetCrons("id,exp,exec_type,exec_target")
	var cronJobs []*CronJob
	var str string
	for _, curCrontab := range newCronJobs {
		cronJobs = append(cronJobs, &CronJob{s, curCrontab})
		str = str + fmt.Sprintf("%d%s%s%s", curCrontab.ID, curCrontab.Exp, curCrontab.ExecType, curCrontab.ExecTarget)
	}
	newCronJobsMD5 := goutils.MD5Str([]byte(str))
	if s.crontabMD5 == newCronJobsMD5 {
		return false
	}
	s.crontabMD5 = newCronJobsMD5
	s.cronJobs = cronJobs
	return true
}

// 计划任务
type CronJob struct {
	scheduler *Scheduler
	*models.Crontab
}

// 运行任务
func (c *CronJob) Run() {
	entry := c.scheduler.c.Entry(cron.EntryID(c.CronId))
	if entry.Prev.Unix() > 0 {
		c.LastExec = entry.Prev
	}
	c.NextExec = entry.Next
	c.Update()

	if c.ExecTarget == "" {
		return
	}
	cronLog := &models.CronLog{
		CID:           c.ID,
		ExecStartTime: time.Now(),
		ExecType:      c.ExecType,
		ExecTarget:    c.ExecTarget,
	}

	if c.ExecType == "shell" {
		cmd := exec.Command(c.ExecTarget)
		err := cmd.Start()
		cronLog.PID = cmd.Process.Pid
		cronLog.IsCrontab = 1
		cronLog.Insert()
		c.scheduler.childProcesses[cronLog.PID] = cmd

		err = cmd.Wait()
		delete(c.scheduler.childProcesses, cronLog.PID)
		if err != nil {
			cronLog.ExecResult = "failed:" + err.Error()
		}
		cronLog.ExecStatus = 1
		cronLog.ExecEndTime = time.Now()
		cronLog.Update()
	}
}

// 当前运行的任务
//type Job struct {
//	Id         int
//	Pid        int
//	Exp        string
//	ExecType   string
//	ExecTarget string
//	StartAt    time.Time
//	EndAt      time.Time
//	Status     int
//	Remark     string
//}

// 运行任务
//func (j *Job) Run() {
//
//}

// 配置
type Config struct {
	Host             string
	Port             int
	User             string
	Password         string
	DBName           string
	CronTableName    string
	CronLogTableName string
	ServicePort      string
}
