package loginRequest

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"strconv"
	"time"
	"u9/tool"
)

//魅族

type meizuChannelRet struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Value   string `json:"message"`
}

type MeiZu struct {
	Lr
	channelRet meizuChannelRet
}

func LrNewMeiZu(channelUserId, token string, args *map[string]interface{}) *MeiZu {
	ret := new(MeiZu)
	ret.Init(channelUserId, token, args)
	return ret
}

func (this *MeiZu) Init(channelUserId, token string, args *map[string]interface{}) {
	this.Lr.Init(channelUserId, token)
	appid := (*args)["MEIZU_APPID"].(string)
	secretKey := (*args)["MEIZU_APPSECRET"].(string)
	this.Method = "POST"
	this.IsHttps = true
	ts := strconv.FormatInt(time.Now().Unix(), 10)
	format := `https://api.game.meizu.com/game/security/checksession?app_id=%s&session_id=%s&uid=%s&ts=%s&sign_type=md5&sign=%s`
	context := "app_id=%s&session_id=%s&ts=%s&uid=%s:%s"
	context = fmt.Sprintf(context, appid, this.token, ts, this.channelUserId, secretKey)
	sign := tool.Md5([]byte(context))
	this.Url = fmt.Sprintf(format, appid, this.token, channelUserId, ts, sign)
}

func (this *MeiZu) ParseChannelRet() (err error) {
	beego.Trace(this.Result)
	err = json.Unmarshal([]byte(this.Result), &this.channelRet)
	return
}

func (this *MeiZu) CheckChannelRet() bool {
	return this.channelRet.Code == 200
}
