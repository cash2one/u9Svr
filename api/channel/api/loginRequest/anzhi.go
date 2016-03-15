package loginRequest

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"strings"
	"time"
)

//安智
//{'msg':'eydjb2RlJzonMTAnfQ==','sc':'10','st':'请求参数错误','time':'20160305112957679'}

type AnZhiChannelRet struct {
	Sc   string `json:"sc"`
	St   string `json:"st"`
	Time string `json:"time"`
	Msg  string `json:"msg"`
}

type AnZhi struct {
	Lr
	channelRet AnZhiChannelRet
}

func LrNewAnZhi(channelUserId, token string, args *map[string]interface{}) *AnZhi {
	ret := new(AnZhi)
	ret.Init(channelUserId, token, args)
	return ret
}

func (this *AnZhi) Init(channelUserId, token string, args *map[string]interface{}) {
	this.Lr.Init(channelUserId, token)
	this.Method = "POST"
	appkey := (*args)["ANZHI_APPKEY"].(string)
	appSecret := (*args)["ANZHI_APPSECRET"].(string)
	var time string = time.Unix(time.Now().Unix(), 0).Format("20060102150405025")
	baseStr := []byte(appkey + this.token + appSecret)
	sign := base64.StdEncoding.EncodeToString(baseStr)

	format := "http://user.anzhi.com/web/api/sdk/third/1/queryislogin?time=%s&appkey=%s&sid=%s&sign=%s"
	this.Url = fmt.Sprintf(format, time, appkey, token, sign)
}

func (this *AnZhi) ParseChannelRet() (err error) {
	beego.Trace(this.Result)
	return json.Unmarshal([]byte(strings.Replace(this.Result, "'", "\"", -1)), &this.channelRet)
}

func (this *AnZhi) CheckChannelRet() bool {
	return this.channelRet.Sc == "1"
}
