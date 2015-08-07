package routers

import (
	"github.com/astaxie/beego"
	"github.com/meinside/rpimongo/controllers"
)

func init() {
	// index
	beego.Router("/", &controllers.IndexController{})

	// links
	beego.Router("/links", &controllers.LinksController{})

	// api
	beego.Router("/api/:method.json", &controllers.ApiController{})
}
