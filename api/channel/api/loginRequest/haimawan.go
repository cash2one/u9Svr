package loginRequest

import (
	"github.com/astaxie/beego"
	"u9/models"
)

//海马玩
type HaiMaWan struct {
	Lr
	args *map[string]interface{}
}

func LrNewHaiMaWan(mlr *models.LoginRequest, args *map[string]interface{}) *HaiMaWan {
	ret := new(HaiMaWan)
	ret.Init(mlr, args)
	return ret
}

func (this *HaiMaWan) Init(mlr *models.LoginRequest, args *map[string]interface{}) {
	this.Lr.Init(mlr)
	this.args = args
	this.Method = "POST"
	this.Url = "http://api.haimawan.com/index.php?m=api&a=validate_token"
}

func (this *HaiMaWan) InitParam() (err error) {
	if err = this.Lr.InitParam(); err != nil {
		return
	}
	appId := (*this.args)["HMKey"].(string)
	this.Req.Param("appid", appId)
	this.Req.Param("t", this.mlr.Token)
	this.Req.Param("uid", this.mlr.ChannelUserid)
	return
}

func (this *HaiMaWan) CheckChannelRet() bool {
	beego.Trace(this.Result)
	return this.Result != "fail"
}
