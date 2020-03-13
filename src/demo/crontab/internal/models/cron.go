package models

import (
	"Demo/crontab/internal/config"
	"fmt"
	"log"
	"time"
)

// 数据库中单行计划任务结构体
type Cron struct {
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

func (c *Cron) Get(selects string) *Cron {
	if selects == "*" {
		selects = "id,title,cron_id,exp,exec_type,exec_target,last_exec,next_exec,is_enable"
	}
	err := DB.Get(c, `select `+selects+` from `+config.CronTableName+` where id=? and is_delete = 0 limit 1`, c.ID)
	if err != nil {
		log.Printf("getCron error: %s\n", err.Error())
	}
	return c
}

func (c *Cron) Update() {
	sql := `update ` + config.CronTableName + ` set cron_id = ?,last_exec =? ,next_exec = ? where id = ?`
	GetDB().MustExec(sql, c.CronId, c.LastExec, c.NextExec, c.ID)
}

func (c *Cron) Enable() {
	sql := `update ` + config.CronTableName + ` set is_enable = ? where id = ?`
	GetDB().MustExec(sql, c.IsEnable, c.ID)
}

func (c *Cron) Del() {
	sql := `update ` + config.CronTableName + ` set is_delete = 1 where id = ?`
	GetDB().MustExec(sql, c.IsEnable, c.ID)
}

func (c *Cron) String() string {
	return fmt.Sprintf("%v", *c)
}

func GetCrons(selects string) []*Cron {
	var jobs []*Cron
	err := DB.Select(&jobs, `select `+selects+` from `+config.CronTableName+` where is_delete = 0`)
	if err != nil {
		log.Printf("getCrons error: %s\n", err.Error())
	}
	return jobs
}
