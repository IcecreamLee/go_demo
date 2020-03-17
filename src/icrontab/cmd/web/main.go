package main

import (
	"github.com/gin-gonic/gin"
	"icrontab/cmd/web/routers"
	"icrontab/internal/config"
)

func main() {
	r := gin.Default()
	routers.SetRouters(r)
	r.Static("/static", "./statics")
	r.StaticFile("/favicon.ico", "./statics/favicon.ico")
	r.LoadHTMLGlob("views/*")
	r.Run(":" + config.WebServicePort) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
