package controllers

import (
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["Website"] = "www.ocean.cri-paris.org/"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "index.tpl"
}
