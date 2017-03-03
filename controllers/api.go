package controllers

import (
	"fmt"
	"strings"

	"github.com/astaxie/beego"

	"github.com/meinside/rpi-tools/status"
)

type ApiResult struct {
	Result string `json:"result"`
	Value  string `json:"value"`
}

// Operations about APIs
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
	if value, err := readValue(method); err == nil {
		res, val = "ok", value
	} else {
		res, val = "error", value
	}

	// enable CORS for dev (swagger)
	if beego.BConfig.RunMode == "dev" {
		c.Ctx.Output.Header("Access-Control-Allow-Origin", "*")
	}

	// output as json
	c.Data["json"] = &ApiResult{
		Result: res,
		Value:  val,
	}
	c.ServeJSON()
}

// Read system values with rpi-tools
func readValue(method string) (result string, err error) {
	switch method {
	case "hostname": // hostname
		return status.Hostname()
	case "uname": // uname -a
		return status.Uname()
	case "uptime": // uptime
		return status.Uptime()
	case "free_spaces": // df -h
		return status.FreeSpaces()
	case "memory_split": // vcgencmd get_mem arm && vcgencmd get_mem gpu
		splits, err := status.MemorySplit()
		return strings.Join(splits, "\n"), err
	case "free_memory": // free -o -h
		return status.FreeMemory()
	case "cpu_temperature": // vcgencmd measure_temp
		return status.CpuTemperature()
	case "cpu_info": //cat /proc/cpuinfo
		return status.CpuInfo()
	default:
		return "Error", fmt.Errorf("No such method: %s", method)
	}
}
