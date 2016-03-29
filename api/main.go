package main

import (
	"github.com/astaxie/beego"
	_ "u9/api/routers"
	_ "u9/models"
	"u9/tool"
)

func main() {
	tool.SetFilelog(true)
	beego.Run()
}
