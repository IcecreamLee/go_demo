package models

import (
	"icrontab/internal/config"
	"icrontab/internal/logger"
	"time"
)

// 计划任务执行日志
type CronLog struct {
	ID            int       `db:"id"`
	CID           int       `db:"cid"`
	PID           int       `db:"pid"`
	IsCrontab     int       `db:"is_crontab"`
	ExecStartTime time.Time `db:"exec_start_time"`
	ExecEndTime   time.Time `db:"exec_end_time"`
	ExecType      string    `db:"exec_type"`
	ExecTarget    string    `db:"exec_target"`
	ExecStatus    int       `db:"exec_status"`
	ExecResult    string    `db:"exec_result"`
}

func (c *CronLog) Get(selects string) *CronLog {
	if selects == "*" {
		selects = "id,cid,pid,is_crontab,exec_start_time,exec_end_time,exec_type,exec_target,exec_status,exec_result"
	}
	err := DB.Get(c, `select `+selects+` from `+config.CronLogTableName+` where id=? limit 1`, c.ID)
	if err != nil {
		logger.Infof("CronLog Get error: %s\n", err.Error())
	}
	return c
}

func (c *CronLog) Insert() int64 {
	sql := `insert into ` + config.CronLogTableName + ` (cid,pid,is_crontab,exec_start_time,exec_type,exec_target,exec_status,exec_result) values (?,?,?,?,?,?,?,?)`
	res := GetDB().MustExec(sql, c.CID, c.PID, c.IsCrontab, c.ExecStartTime, c.ExecType, c.ExecTarget, c.ExecStatus, c.ExecResult)
	id, err := res.LastInsertId()
	if err != nil {
		return 0
	}
	c.ID = int(id)
	return id
}

func (c *CronLog) Update() {
	sql := `update ` + config.CronLogTableName + ` set pid=?,exec_end_time=?,exec_status=?,exec_result=? where id=?`
	GetDB().MustExec(sql, c.PID, c.ExecEndTime, c.ExecStatus, c.ExecResult, c.ID)
}

func GetCronLogsByCID(cid int, selects string) []*CronLog {
	var logs []*CronLog
	err := DB.Select(&logs, `select `+selects+` from `+config.CronLogTableName+` where cid=? order by id desc limit 100`, cid)
	if err != nil {
		logger.Infof("GetCronLogsByCID error: %s\n", err.Error())
	}
	return logs
}
