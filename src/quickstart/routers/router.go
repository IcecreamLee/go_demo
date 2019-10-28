package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"quickstart/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/login", &controllers.LoginController{})
	beego.Router("/user", &controllers.UserController{})
	beego.Router("/id", &controllers.IDController{})
	beego.Get("/hello", func(ctx *context.Context) {
		ctx.Output.Body([]byte("hello world!"))
	})
}
