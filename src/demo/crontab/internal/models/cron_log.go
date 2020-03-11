package models

import (
	"Demo/crontab/internal/config"
	"time"
)

// 计划任务执行日志
type CronLog struct {
	ID            int       `db:"id"`
	CID           int       `db:"cid"`
	PID           int       `db:"pid"`
	ExecStartTime time.Time `db:"exec_start_time"`
	ExecEndTime   time.Time `db:"exec_end_time"`
	ExecType      string    `db:"exec_type"`
	ExecTarget    string    `db:"exec_target"`
	ExecStatus    int       `db:"exec_status"`
	ExecResult    string    `db:"exec_result"`
}

func (cl *CronLog) Insert() int64 {
	sql := `insert into ` + config.CronLogTableName + ` (cid,pid,exec_start_time,exec_type,exec_target,exec_status,exec_result) values (?,?,?,?,?,?,?)`
	res := GetDB().MustExec(sql, cl.CID, cl.PID, cl.ExecStartTime, cl.ExecType, cl.ExecTarget, cl.ExecStatus, cl.ExecResult)
	id, err := res.LastInsertId()
	if err != nil {
		return 0
	}
	cl.ID = int(id)
	return id
}

func (cl *CronLog) Update() {
	sql := `update ` + config.CronLogTableName + ` set pid=?,exec_end_time=?,exec_status=?,exec_result=? where id=?`
	GetDB().MustExec(sql, cl.PID, cl.ExecEndTime, cl.ExecStatus, cl.ExecResult, cl.ID)
}
