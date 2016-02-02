package controllers

import (
	"github.com/astaxie/beego"
)

type LinksController struct {
	beego.Controller
}

func (c *LinksController) Get() {
	c.Layout = "layouts/layout.html"
	c.TplName = "links.tpl"
}
