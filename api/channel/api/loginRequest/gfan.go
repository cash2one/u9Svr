package loginRequest

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
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

func LrNewGFan(channelUserId, token string, args *map[string]interface{}) *GFan {
	ret := new(GFan)
	ret.Init(channelUserId, token, args)
	return ret
}

func (this *GFan) Init(channelUserId, token string, args *map[string]interface{}) {
	this.Lr.Init(channelUserId, token)
	format := "http://api.gfan.com/uc1/common/verify_token?token=%s"
	this.Url = fmt.Sprintf(format, token)
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
