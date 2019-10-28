package routers

import (
	"framework/controllers"
	"framework/icecream"
)

func init() {
	icecream.AddRouter("/index", &controllers.IndexController{})
}
