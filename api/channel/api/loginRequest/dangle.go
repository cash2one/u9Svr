package loginRequest

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"u9/tool"
)

//当乐

type DangleChannelRet struct {
	Valid    string `json:"valid"`
	Roll     bool   `json:"roll"`
	Interval int    `json:"interval"`
	Times    int    `json:"times"`
	Code     int    `json:"msg_code"`
	Desc     string `json:"msg_desc"`
}

type Dangle struct {
	Lr
	channelRet DangleChannelRet
}

func LrNewDangle(channelUserId, token string, args *map[string]interface{}) *Dangle {
	ret := new(Dangle)
	ret.Init(channelUserId, token, args)
	return ret
}

func (this *Dangle) Init(channelUserId, token string, args *map[string]interface{}) {
	this.Lr.Init(channelUserId, token)
	appid := (*args)["DANGLE_SDK_APPID"].(string)
	appkey := (*args)["DANGLE_SDK_APPKEY"].(string)
	sign := tool.Md5([]byte(appid + "|" + appkey + "|" + token + "|" + channelUserId))
	format := "http://ngsdk.d.cn/api/cp/checkToken?appid=%s&umid=%s&token=%s&sig=%s"
	this.Url = fmt.Sprintf(format, appid, channelUserId, token, sign)
}

func (this *Dangle) ParseChannelRet() (err error) {
	beego.Trace(this.Result)
	return json.Unmarshal([]byte(this.Result), &this.channelRet)
}

func (this *Dangle) CheckChannelRet() bool {
	return this.channelRet.Code == 2000
}
