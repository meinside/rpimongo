package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["github.com/meinside/rpimongo/controllers:ApiController"] = append(beego.GlobalControllerRouter["github.com/meinside/rpimongo/controllers:ApiController"],
		beego.ControllerComments{
			"Get",
			`/:method`,
			[]string{"get"},
			nil})

}
