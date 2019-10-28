package controllers

import (
	"github.com/astaxie/beego"
)

type LoginController struct {
	beego.Controller
}

func (this *LoginController) Get() {
	//fmt.Fprintf(this.Ctx.ResponseWriter, this.Ctx.Request.URL.Query()["a"][0])
	this.Data["json"] = this.Ctx.Request.URL.Query()
	this.ServeJSON()
	//this.Ctx.Output.Body([]byte("xxx"))
	//this.TplName = "login.html"
}

