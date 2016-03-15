package cp

import (
	//"github.com/astaxie/beego"
	"os"
	"runtime"
	"u9/models"
)

type IndexController struct {
	BaseController
}

func (this *IndexController) Index() {
	this.updateData()
	this.display()
}

func (this *IndexController) updateData() {
	this.BaseController.updateData()
	this.Data["hostname"], _ = os.Hostname()
	this.Data["gover"] = runtime.Version()
	this.Data["os"] = runtime.GOOS
	this.Data["cpunum"] = runtime.NumCPU()
	this.Data["arch"] = runtime.GOARCH
	this.Data["usernum"], _ = new(models.Cp).Query().Count()
}
