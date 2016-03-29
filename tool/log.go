package tool

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

func SetFilelog(active bool) {
	if active {
		beego.BeeLogger.SetLogger("file", `{"filename":"log/u9.log", "level":7}`)
		beego.SetLogFuncCall(true)
		beego.BeeLogger.Async()
	} else {
		//beego.BeeLogger.SetLogger("console", `{"level":7}`)
	}
	beego.BConfig.Log.AccessLogs = active
	beego.BConfig.Log.FileLineNum = active
	orm.Debug = active
}
