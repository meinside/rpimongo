package controllers

import (
	"github.com/astaxie/beego"
)

type IndexController struct {
	beego.Controller
}

func (c *IndexController) Get() {
	c.Layout = "layouts/layout.html"
	c.TplNames = "index.tpl"
}
