package routers

import (
	"github.com/astaxie/beego"
	"u9/www/controllers/cp"
	"u9/www/controllers/manager"
)

func init() {
	//manager
	beego.Router("/manager", &manager.IndexController{}, "*:Index")
	beego.Router("/manager/login", &manager.LoginController{}, "*:Login")
	beego.Router("/manager/logout", &manager.BaseController{}, "*:Logout")
	beego.Router("/manager/profile", &manager.ProfileController{}, "*:Profile")

	beego.Router("/manager/cp/list", &manager.CpController{}, "*:List")
	beego.Router("/manager/cp/add", &manager.CpController{}, "*:Add")
	beego.Router("/manager/cp/edit", &manager.CpController{}, "*:Edit")
	beego.Router("/manager/cp/delete", &manager.CpController{}, "*:Delete")

	beego.Router("/manager/channel/list", &manager.ChannelController{}, "*:List")
	beego.Router("/manager/channel/add", &manager.ChannelController{}, "*:Add")
	beego.Router("/manager/channel/edit", &manager.ChannelController{}, "*:Edit")
	beego.Router("/manager/channel/delete", &manager.ChannelController{}, "*:Delete")

	beego.Router("/manager/statistic/payList",&manager.StatisticController{},"*:PayList")

	//cp
	beego.Router("/cp", &cp.IndexController{}, "*:Index")
	beego.Router("/cp/login", &cp.LoginController{}, "*:Login")
	beego.Router("/cp/logout", &cp.BaseController{}, "*:Logout")
	beego.Router("/cp/profile", &cp.ProfileController{}, "*:Profile")

	beego.Router("/cp/product/list", &cp.ProductController{}, "*:List")
	beego.Router("/cp/product/edit", &cp.ProductController{}, "*:Edit")
	beego.Router("/cp/product/delete", &cp.ProductController{}, "*:Delete")

	beego.Router("/cp/package/list", &cp.PackageController{}, "*:List")
	beego.Router("/cp/package/add", &cp.PackageController{}, "*:Add")
	beego.Router("/cp/package/getVcParam", &cp.PackageController{}, "*:GetVerAndClByPid")
	beego.Router("/cp/package/package", &cp.PackageController{}, "*:Package")
	beego.Router("/cp/package/download", &cp.PackageController{}, "*:Download")
	beego.Router("/cp/package/delete", &cp.PackageController{}, "*:Delete")

	
}
