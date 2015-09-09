package main

import (
	"github.com/astaxie/beego"
	_ "github.com/meinside/rpimongo/docs"
	_ "github.com/meinside/rpimongo/routers"
)

func main() {
	if beego.RunMode == "dev" {
		beego.DirectoryIndex = true
		beego.StaticDir["/swagger"] = "swagger"
	}

	beego.Run()
}
