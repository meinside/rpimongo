// @APIVersion 0.0.1
// @Title RPiMonGo Test API
// @Description For testing JSON APIs for RPiMonGo
// @Contact meinside@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
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
