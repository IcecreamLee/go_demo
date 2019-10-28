package controllers

import (
	"fmt"
	"framework/icecream/core"
)

type IndexController struct {
	core.Controller
}

func (c *IndexController) Get() {
	fmt.Print("xxxxxxxxxxx")
	//fmt.Fprintf(c.Request.Response, "URL%s", c.Request.Request)
}
