package controllers

import (
	"github.com/astaxie/beego"
	"github.com/meinside/rpimongo/lib"
)

type ApiResult struct {
	Result string `json:"result"`
	Value  string `json:"value"`
}

type ApiController struct {
	beego.Controller
}

// @Title read
// @Description Query status from server
// @Param method path string   true        "status method"
// @Success 200 {string} json
// @router /:method [get]
func (c *ApiController) Get() {
	method := c.Ctx.Input.Param(":method")

	var res, val string
	if value, success := rpi.ReadValue(method); success {
		res, val = "ok", value
	} else {
		res, val = "error", value
	}

	// output as json
	c.Data["json"] = &ApiResult{
		Result: res,
		Value:  val,
	}
	c.ServeJson()
}
