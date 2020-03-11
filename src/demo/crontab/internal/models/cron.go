package models

import (
	"Demo/crontab/internal/config"
	"fmt"
	"time"
)

// 数据库中单行计划任务结构体
type Cron struct {
	ID         int       `db:"id"`
	CronId     int       `db:"cron_id"`
	Exp        string    `db:"exp"`
	ExecType   string    `db:"exec_type"`
	ExecTarget string    `db:"exec_target"`
	LastExec   time.Time `db:"last_exec"`
	NextExec   time.Time `db:"next_exec"`
}

func (cronJob *Cron) Update() {
	sql := `update ` + config.CronTableName + ` set cron_id = ?,last_exec =? ,next_exec = ? where id = ?`
	GetDB().MustExec(sql, cronJob.CronId, cronJob.LastExec, cronJob.NextExec, cronJob.ID)
}

func (cronJob *Cron) String() string {
	return fmt.Sprintf("%v", *cronJob)
}
