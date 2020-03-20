package main

import (
	"github.com/IcecreamLee/goutils"
	"github.com/gin-gonic/gin"
	"html/template"
	"icrontab/cmd/web/routers"
	"icrontab/internal/config"
)

func main() {
	r := gin.Default()
	routers.SetRouters(r)

	r.SetFuncMap(template.FuncMap{
		"formatDate":     goutils.DateFormat,
		"formatDatetime": goutils.DatetimeFormat,
	})

	r.Static("/static", "./statics")
	r.StaticFile("/favicon.ico", "./statics/favicon.ico")
	r.LoadHTMLGlob("views/*")

	r.Run(":" + config.WebServicePort) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
