package models

import (
	"fmt"
	"icrontab/internal/config"
	"icrontab/internal/logger"
	"time"
)

// 数据库中单行计划任务结构体
type Crontab struct {
	ID         int       `db:"id"`
	Title      string    `db:"title"`
	CronId     int       `db:"cron_id"`
	Exp        string    `db:"exp"`
	ExecType   string    `db:"exec_type"`
	ExecTarget string    `db:"exec_target"`
	LastExec   time.Time `db:"last_exec"`
	NextExec   time.Time `db:"next_exec"`
	IsEnable   int       `db:"is_enable"`
}

func (c *Crontab) Get(selects string) *Crontab {
	if selects == "*" {
		selects = "id,title,cron_id,exp,exec_type,exec_target,last_exec,next_exec,is_enable"
	}
	err := DB.Get(c, `select `+selects+` from `+config.CronTableName+` where id=? and is_delete = 0 limit 1`, c.ID)
	if err != nil {
		logger.Infof("Crontab Get error: %s\n", err.Error())
	}
	return c
}

func (c *Crontab) Update() {
	sql := `update ` + config.CronTableName + ` set cron_id = ?,last_exec =? ,next_exec = ? where id = ?`
	GetDB().MustExec(sql, c.CronId, c.LastExec, c.NextExec, c.ID)
}

func (c *Crontab) Enable() {
	sql := `update ` + config.CronTableName + ` set is_enable = ? where id = ?`
	GetDB().MustExec(sql, c.IsEnable, c.ID)
}

func (c *Crontab) Del() {
	sql := `update ` + config.CronTableName + ` set is_delete = 1 where id = ?`
	GetDB().MustExec(sql, c.IsEnable, c.ID)
}

func (c *Crontab) String() string {
	return fmt.Sprintf("%v", *c)
}

func GetCrons(selects string) []*Crontab {
	var jobs []*Crontab
	err := DB.Select(&jobs, `select `+selects+` from `+config.CronTableName+` where is_delete = 0`)
	if err != nil {
		logger.Infof("GetCrons error: %s\n", err.Error())
	}
	return jobs
}

func GetEnabledCrons(selects string) []*Crontab {
	var jobs []*Crontab
	err := DB.Select(&jobs, `select `+selects+` from `+config.CronTableName+` where is_delete=0 and is_enable=1`)
	if err != nil {
		logger.Infof("GetCrons error: %s\n", err.Error())
	}
	return jobs
}
