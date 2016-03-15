package main

import (
	"github.com/astaxie/beego"
	_ "u9/models"
	_ "u9/www/routers"
)

func main() {
	beego.Run()
}
