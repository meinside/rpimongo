package routers

import (
	"github.com/meinside/rpimongo/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})
}
