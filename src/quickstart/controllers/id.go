package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"quickstart/lib"
)

type IDController struct {
	beego.Controller
	ids map[int64]int64
}

func (c *IDController) Get() {
	// 另一种map的声明方式
	c.ids = make(map[int64]int64)

	//fmt.Fprintf(c.Ctx.ResponseWriter, "Response:\n")
	//timestamp := time.Now().Unix()
	for i := 0; i < 10000; i++ {
		//	go func() {
		c.getID()
		//fmt.Fprintf(c.Ctx.ResponseWriter, "id%4d: %d\n", i, id)
		//}()
	}
	b, _ := json.Marshal(c.ids)
	fmt.Fprintf(c.Ctx.ResponseWriter, "%s", b)
	//c.Data["User"] = time.Now().Unix() - timestamp
	//c.TplName = "user.html"
}

func (c *IDController) getID() {
	id := lib.GetIDInstance().Get()
	c.ids[id] = id
}
