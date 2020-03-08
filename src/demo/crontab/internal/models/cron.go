package models

import "C"
import (
	"Demo/crontab/internal/crontab"
	"fmt"
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

func (cronJob *Cron) Update() {
	sql := `update ` + crontab.Conf.CronTableName + ` set cron_id = ?,last_exec =? ,next_exec = ? where id = ?`
	GetDB().MustExec(sql, cronJob.CronId, cronJob.LastExec, cronJob.NextExec, cronJob.ID)
}

func (cronJob *Cron) String() string {
	return fmt.Sprintf("%v", *cronJob)
}
