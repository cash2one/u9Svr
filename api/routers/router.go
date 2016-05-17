package routers

import (
	"github.com/astaxie/beego"
	"u9/api/controllers/login"
	"u9/api/controllers/pay"
	"u9/api/controllers/test"
)

func init() {
	//login api
	beego.Router("/api/gameLoginRequest", &login.LoginController{}, "*:LoginRequest")
	beego.Router("/api/validateGameLogin", &login.LoginController{}, "*:ValidateLogin")

	//pay api
	beego.Router("/api/gamePayRequest", &pay.PayController{}, "*:PayRequest")
	beego.Router("/api/channelPayNotify/?:productId/?:channelId", &pay.PayController{}, "*:ChannelPayNotify")

	//test
	beego.Router("/test", &test.Test{}, "*:Test1")
}
