package controllers

import (
	"bytes"
	"github.com/IcecreamLee/goutils"
	"github.com/gin-gonic/gin"
	"icrontab/cmd/web/helpers"
	"icrontab/internal/config"
	"icrontab/internal/models"
	"os/exec"
	"time"
)

var childProcesses = make(map[int]*exec.Cmd)

// 顶层iframe页
func Login(c *gin.Context) {
	session := helpers.GetSession(c)
	session.Values = make(map[interface{}]interface{})

	username := c.PostForm("username")
	password := c.PostForm("password")
	if username == config.WebUserName && password == config.WebPassword {
		session.Values["username"] = username
		session.Values["password"] = password
		_ = session.Save(c.Request, c.Writer)
		c.JSON(200, gin.H{"code": 0, "msg": "登录成功"})
	} else if username != "" || password != "" {
		_ = session.Save(c.Request, c.Writer)
		c.JSON(200, gin.H{"code": 1, "msg": "登录失败"})
	} else {
		_ = session.Save(c.Request, c.Writer)
		c.HTML(200, "login.html", gin.H{
			"title": "登录",
		})
	}
}

func CheckLogin(c *gin.Context) {
	if !isLogin(c) {
		c.Redirect(307, "/login")
		return
	}
	c.Next()
}

func CheckLoginForAjax(c *gin.Context) {
	if !isLogin(c) {
		c.JSON(200, gin.H{"code": 1, "msg": "登录失效"})
		return
	}
	c.Next()
}

// 顶层iframe页
func Index(c *gin.Context) {
	c.HTML(200, "index.html", gin.H{
		"title": "首页",
	})
}

// 任务列表页
func Crons(c *gin.Context) {
	crons := models.GetCrons("id,title,exp,exec_type,exec_target,last_exec,next_exec,is_enable")
	println(crons)
	c.HTML(200, "crons.html", gin.H{
		"title": "任务管理",
		"crons": crons,
	})
}

// 任务历史列表页
func Logs(c *gin.Context) {
	cid := goutils.ToInt(c.Query("cid"))
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	logs := models.GetCronLogsByCID(cid, "id,cid,pid,is_crontab,exec_start_time,exec_end_time,exec_type,exec_target,exec_status,exec_result")
	c.HTML(200, "logs.html", gin.H{
		"title":     "任务历史",
		"logs":      logs,
		"cid":       cid,
		"startDate": startDate,
		"endDate":   endDate,
	})
}

// 新增任务
func Add(c *gin.Context) {
	c.HTML(200, "edit.html", gin.H{
		"title":   "新增任务",
		"crontab": &models.Crontab{},
	})
}

// 编辑任务
func Edit(c *gin.Context) {
	id := c.Query("id")
	crontab := (&models.Crontab{ID: goutils.ToInt(id)}).Get("*")
	c.HTML(200, "edit.html", gin.H{
		"title":   "编辑任务",
		"crontab": crontab,
	})
}

func SaveCrontab(c *gin.Context) {
	id := goutils.ToInt(c.PostForm("id"))
	crontab := &models.Crontab{}
	data := map[string]interface{}{
		"title":       c.PostForm("title"),
		"exp":         c.PostForm("exp"),
		"exec_type":   c.PostForm("exec_type"),
		"exec_target": c.PostForm("exec_target"),
		"is_enable":   goutils.ToInt(c.PostForm("is_enable")),
	}
	if id > 0 {
		err := models.UpdateWithMap(
			crontab.GetTableName(),
			data,
			map[string]interface{}{
				"id": id,
			},
		)
		if err != nil {
			c.JSON(200, gin.H{"code": 1, "msg": "更新失败"})
		} else {
			c.JSON(200, gin.H{"code": 1, "msg": "更新成功"})
		}
	} else {
		err := models.InsertWithMap(crontab.GetTableName(), data)
		if err != nil {
			c.JSON(200, gin.H{"code": 1, "msg": "新增失败"})
		} else {
			c.JSON(200, gin.H{"code": 1, "msg": "新增成功"})
		}
	}
}

// 删除任务
func Del(c *gin.Context) {
	cid := c.PostForm("cid")
	(&models.Crontab{ID: goutils.ToInt(cid)}).Del()
	c.JSON(200, gin.H{"code": 0, "msg": "操作成功"})
}

// 启用/禁用任务
func Enable(c *gin.Context) {
	cid := c.PostForm("cid")
	enable := c.PostForm("enable")
	if enable != "1" {
		enable = "0"
	}
	(&models.Crontab{ID: goutils.ToInt(cid), IsEnable: goutils.ToInt(enable)}).Enable()
	c.JSON(200, gin.H{"code": 0, "msg": "操作成功"})
}

// 运行任务
func Run(c *gin.Context) {
	cid := c.PostForm("cid")
	cron := (&models.Crontab{ID: goutils.ToInt(cid)}).Get("*")
	if cron.ExecTarget == "" {
		c.JSON(200, gin.H{"code": 0, "msg": "操作成功"})
		return
	}

	cronLog := &models.CronLog{
		CID:           cron.ID,
		ExecStartTime: time.Now(),
		ExecType:      cron.ExecType,
		ExecTarget:    cron.ExecTarget,
	}

	if cron.ExecType == "shell" {
		cmd := exec.Command(cron.ExecTarget)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		err := cmd.Start()
		cronLog.PID = cmd.Process.Pid
		cronLog.Insert()
		childProcesses[cronLog.PID] = cmd

		err = cmd.Wait()
		delete(childProcesses, cronLog.PID)
		if err == nil {
			cronLog.ExecResult = "finished:"
			c.JSON(200, gin.H{"code": 0, "msg": "运行完成"})
		} else {
			cronLog.ExecResult = "failed:" + err.Error()
			c.JSON(200, gin.H{"code": 1, "msg": "运行失败：" + err.Error()})
		}
		cronLog.ExecStatus = 1
		cronLog.ExecEndTime = time.Now()
		cronLog.ExecResult += goutils.SubStr(out.String(), 0, 500)
		cronLog.Update()
	}
}

// 停止任务
func Stop(c *gin.Context) {
	id := c.PostForm("id")
	cronLog := (&models.CronLog{ID: goutils.ToInt(id)}).Get("*")

	if cronLog.IsCrontab == 1 {
		res := goutils.HttpPost("http://localhost:"+config.ServicePort+"/stop", "id="+id, goutils.HttpContentTypeForm)
		c.Header("Content-Type", goutils.HttpContentTypeJson)
		c.String(200, res)
		return
	}

	var err error
	cmd, ok := childProcesses[cronLog.PID]
	if ok {
		err = cmd.Process.Kill()
		delete(childProcesses, cronLog.PID)
		if err != nil {
			c.JSON(200, gin.H{"code": 1, "msg": "操作失败，" + err.Error()})
			return
		}
	} else {
		c.JSON(200, gin.H{"code": 1, "msg": "操作失败，子进程不存在"})
		return
	}
	c.JSON(200, gin.H{"code": 0, "msg": "操作成功"})
}

// 停止任务
func Restart(c *gin.Context) {
	res := goutils.HttpPost("http://localhost:"+config.ServicePort+"/restart", "", goutils.HttpContentTypeForm)
	c.Header("Content-Type", goutils.HttpContentTypeJson)
	c.String(200, res)
	return
}

func isLogin(c *gin.Context) bool {
	session := helpers.GetSession(c)
	if session.Values["username"] == config.WebUserName && session.Values["password"] == config.WebPassword {
		return true
	}
	return false
}
