package models

import (
	"Demo/crontab/internal/crontab"
	"time"
)

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

func (cl *CronLog) Insert() {
	sql := `Insert into ` + crontab.Conf.CronLogTableName + ` (cid,exec_start_time,exec_end_time,exec_type,exec_target,exec_status,exec_result) values (?,?,?,?,?,?,?)`
	GetDB().MustExec(sql, cl.CID, cl.ExecStartTime, cl.ExecEndTime, cl.ExecType, cl.ExecTarget, cl.ExecStatus, cl.ExecResult)
}
