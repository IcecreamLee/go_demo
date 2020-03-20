package routers

import (
	"github.com/gin-gonic/gin"
	"icrontab/cmd/web/controllers"
)

func SetRouters(r *gin.Engine) {
	r.GET("/login", controllers.Login)
	r.POST("/login", controllers.Login)
	r.GET("/", controllers.CheckLogin, controllers.Index)
	r.GET("/index", controllers.CheckLogin, controllers.Index)
	r.GET("/crons", controllers.CheckLogin, controllers.Crons)
	r.GET("/logs", controllers.CheckLogin, controllers.Logs)
	r.GET("/add", controllers.CheckLogin, controllers.Add)
	r.POST("/add", controllers.CheckLoginForAjax, controllers.SaveCrontab)
	r.GET("/edit", controllers.CheckLogin, controllers.Edit)
	r.POST("/edit", controllers.CheckLoginForAjax, controllers.SaveCrontab)
	r.POST("/del", controllers.CheckLoginForAjax, controllers.Del)
	r.POST("/enable", controllers.CheckLoginForAjax, controllers.Enable)
	r.POST("/run", controllers.CheckLoginForAjax, controllers.Run)
	r.POST("/stop", controllers.CheckLoginForAjax, controllers.Stop)
	r.POST("/restart", controllers.CheckLoginForAjax, controllers.Restart)
}
