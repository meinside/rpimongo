package routers

import (
	"github.com/astaxie/beego"
	"github.com/meinside/rpimongo/controllers"
	"github.com/yvasiyarov/beego_gorelic"
)

func init() {
	// newrelic monitoring
	beego_gorelic.InitNewrelicAgent()

	// index
	beego.Router("/", &controllers.IndexController{})

	// links
	beego.Router("/links", &controllers.LinksController{})

	// api
	beego.Router("/api/:method.json", &controllers.ApiController{})
}
