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

	// api (namespace for swagger)
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/api",
			beego.NSInclude(
				&controllers.ApiController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
