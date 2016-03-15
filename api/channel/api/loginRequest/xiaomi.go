package loginRequest

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"u9/tool"
)

type XiaomiChannelRet struct {
	Code int    `json:"errcode"`
	Desc string `json:"errMsg"`
}

type Xiaomi struct {
	Lr
	channelRet XiaomiChannelRet
}

func LrNewXiaomi(channelUserId, token string, args *map[string]interface{}) *Xiaomi {
	ret := new(Xiaomi)
	ret.Init(channelUserId, token, args)
	return ret
}

func (this *Xiaomi) Init(channelUserId, token string, args *map[string]interface{}) {
	this.Lr.Init(channelUserId, token)
	appid := (*args)["XIAOMI_APPID"].(string)
	secretkey := (*args)["XIAOMI_SECRETKEY"].(string)
	signContext := fmt.Sprintf("appId=%s&session=%s&uid=%s", appid, token, channelUserId)
	sign := fmt.Sprintf("%x", string(tool.HmacSHA1Encrypt(signContext, secretkey)))
	format := "http://mis.migc.xiaomi.com/api/biz/service/verifySession.do?" + signContext + "&signature=%s"
	this.Url = fmt.Sprintf(format, sign)
	beego.Trace(this.Url)
}

func (this *Xiaomi) ParseChannelRet() (err error) {
	beego.Trace(this.Result)
	return json.Unmarshal([]byte(this.Result), &this.channelRet)
}

func (this *Xiaomi) CheckChannelRet() bool {
	return this.channelRet.Code == 200
}
