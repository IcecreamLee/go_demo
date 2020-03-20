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

func (c *Crontab) GetTableName() string {
	return config.CronTableName
}

func (c *Crontab) Get(selects string) *Crontab {
	if selects == "*" {
		selects = "id,title,cron_id,exp,exec_type,exec_target,last_exec,next_exec,is_enable"
	}
	err := DB.Get(c, `select `+selects+` from `+c.GetTableName()+` where id=? and is_delete = 0 limit 1`, c.ID)
	logger.IfError("Failed to get crontab: %s", err)
	return c
}

func (c *Crontab) Update() {
	sql := `update ` + c.GetTableName() + ` set cron_id = ?,last_exec =? ,next_exec = ? where id = ?`
	GetDB().MustExec(sql, c.CronId, c.LastExec, c.NextExec, c.ID)
}

func (c *Crontab) Enable() {
	sql := `update ` + c.GetTableName() + ` set is_enable = ? where id = ?`
	GetDB().MustExec(sql, c.IsEnable, c.ID)
}

func (c *Crontab) Del() {
	sql := `update ` + c.GetTableName() + ` set is_delete = 1 where id = ?`
	GetDB().MustExec(sql, c.ID)
}

func (c *Crontab) String() string {
	return fmt.Sprintf("%v", *c)
}

func GetCrons(selects string) []*Crontab {
	var jobs []*Crontab
	err := DB.Select(&jobs, `select `+selects+` from `+config.CronTableName+` where is_delete = 0`)
	logger.IfError("Failed to get crontabs: %s", err)
	return jobs
}

func GetEnabledCrons(selects string) []*Crontab {
	var jobs []*Crontab
	err := DB.Select(&jobs, `select `+selects+` from `+config.CronTableName+` where is_delete=0 and is_enable=1`)
	logger.IfError("Failed to get enabled crontabs: %s", err)
	return jobs
}
