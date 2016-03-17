package loginRequest

import (
	"github.com/astaxie/beego"
)

//海马玩
type HaiMaWan struct {
	Lr
	args *map[string]interface{}
}

func LrNewHaiMaWan(channelUserId, token string, args *map[string]interface{}) *HaiMaWan {
	ret := new(HaiMaWan)
	ret.Init(channelUserId, token, args)
	return ret
}

func (this *HaiMaWan) Init(channelUserId, token string, args *map[string]interface{}) {
	this.Lr.Init(channelUserId, token)
	this.args = args
	this.Method = "POST"
	this.Url = "http://api.haimawan.com/index.php?m=api&a=validate_token"
}

func (this *HaiMaWan) InitParam() {
	this.Lr.InitParam()
	appId := (*this.args)["HMKey"].(string)
	this.Req.Param("appid", appId)
	this.Req.Param("t", this.token)
	this.Req.Param("uid", this.channelUserId)
}

func (this *HaiMaWan) CheckChannelRet() bool {
	beego.Trace(this.Result)
	return this.Result != "fail"
}
