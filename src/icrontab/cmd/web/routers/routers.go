package routers

import (
	"github.com/gin-gonic/gin"
	"icrontab/cmd/web/controllers"
)

func SetRouters(r *gin.Engine) {
	r.GET("/", controllers.CheckLogin, controllers.Index)
	r.GET("/index", controllers.CheckLogin, controllers.Index)
	r.GET("/login", controllers.Login)
	r.POST("/login", controllers.Login)
	r.GET("/crons", controllers.CheckLogin, controllers.Crons)
	r.GET("/logs", controllers.CheckLogin, controllers.Logs)
	r.GET("/add", controllers.Add)
	r.GET("/edit", controllers.Edit)
	r.POST("/del", controllers.Del)
	r.POST("/enable", controllers.Enable)
	r.GET("/run", controllers.Run)
	r.POST("/run", controllers.Run)
	r.GET("/stop", controllers.Stop)
	r.POST("/stop", controllers.Stop)
	r.POST("/restart", controllers.Restart)
}
