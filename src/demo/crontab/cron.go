package main

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"os/exec"
	"time"
)

// 数据库中单行计划任务结构体
type Cron struct {
	ID         int
	CronId     int `db:"cron_id"`
	Exp        string
	ExecType   string    `db:"exec_type"`
	ExecTarget string    `db:"exec_target"`
	LastExec   time.Time `db:"last_exec"`
	NextExec   time.Time `db:"next_exec"`
	IsDelete   int       `db:"is_delete"`
	AddAt      time.Time `db:"add_at"`
}

// 计划任务运行
func (cronJob *Cron) Run() {
	entry := c.goCron.Entry(cron.EntryID(cronJob.CronId))

	if entry.Prev.Unix() > 0 {
		cronJob.LastExec = entry.Prev
	}
	cronJob.NextExec = entry.Next
	cronJob.Update()

	if cronJob.ExecTarget == "" {
		return
	}

	now := time.Now()
	cronLog := CronLog{cronJob.ID, now, now, cronJob.ExecType, cronJob.ExecTarget, 0, ""}

	if cronJob.ExecType == "shell" {
		cmd := exec.Command(cronJob.ExecTarget)
		_, err := cmd.CombinedOutput()
		if err == nil {
			cronLog.ExecStatus = 1
		} else {
			logger.Printf("[error] cronJob run failed:%s\n", err)
			cronLog.ExecResult = "failed:" + err.Error()
		}
		cronLog.ExecEndTime = time.Now()
		cronLog.insert()
	}
}

func (cronJob *Cron) Update() {
	sql := `update ` + config.cronTableName + ` set cron_id = ?,last_exec =? ,next_exec = ? where id = ?`
	db.MustExec(sql, cronJob.CronId, cronJob.LastExec, cronJob.NextExec, cronJob.ID)
}

func (cronJob *Cron) String() string {
	return fmt.Sprintf("%v", *cronJob)
}

// 计划任务执行日志
type CronLog struct {
	CID           int       `db:"cid"`
	ExecStartTime time.Time `db:"exec_start_time"`
	ExecEndTime   time.Time `db:"exec_end_time"`
	ExecType      string    `db:"exec_type"`
	ExecTarget    string    `db:"exec_target"`
	ExecStatus    int       `db:"exec_status"`
	ExecResult    string    `db:"exec_result"`
}

func (cl *CronLog) insert() {
	sql := `insert into ` + config.cronLogTableName + ` (cid,exec_start_time,exec_end_time,exec_type,exec_target,exec_status,exec_result) values (?,?,?,?,?,?,?)`
	db.MustExec(sql, cl.CID, cl.ExecStartTime, cl.ExecEndTime, cl.ExecType, cl.ExecTarget, cl.ExecStatus, cl.ExecResult)
}
