package main

import (
	"github.com/astaxie/beego"
	_ "u9/models"
	"u9/tool"
	_ "u9/www/routers"
)

func main() {
	tool.SetFilelog(true)
	beego.Run()
}
