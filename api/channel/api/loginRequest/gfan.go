package loginRequest

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"u9/models"
)

//机锋网

type GFanChannelRet struct {
	ResultCode int    `json:"resultCode"`
	UId        uint64 `json:"uid"`
}

type GFan struct {
	Lr
	channelRet GFanChannelRet
	args       *map[string]interface{}
}

func LrNewGFan(mlr *models.LoginRequest, args *map[string]interface{}) *GFan {
	ret := new(GFan)
	ret.Init(mlr, args)
	return ret
}

func (this *GFan) Init(mlr *models.LoginRequest, args *map[string]interface{}) {
	this.Lr.Init(mlr)
	format := "http://api.gfan.com/uc1/common/verify_token?token=%s"
	this.Url = fmt.Sprintf(format, this.mlr.Token)
	this.args = args
}

func (this *GFan) InitParam() {
	this.Lr.InitParam()
	uid := (*this.args)["GFAN_UID"].(string)
	this.Req.Header("channelID", uid)
}

func (this *GFan) ParseChannelRet() (err error) {
	beego.Trace(this.Result)
	err = json.Unmarshal([]byte(this.Result), &this.channelRet)
	return
}

func (this *GFan) CheckChannelRet() bool {
	return this.channelRet.ResultCode == 1
}
