package main

import (
	"github.com/astaxie/beego"
	_ "u9/api/routers"
	_ "u9/models"
)

func main() {
	beego.Run()
	// beego.SetLevel(beego.LevelDebug)
	// beego.SetLogFuncCall(true)
	// beego.SetLogger("file", `{"filename":"api.log"}`)
}
