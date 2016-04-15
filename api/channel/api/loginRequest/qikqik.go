package loginRequest

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"u9/tool"
)

//7k7k

type QikQikChannelRet struct {
	Status    int `json:"status"`
	Msg     string   `json:"msg"`
}

type QikQik struct {
	Lr
	channelRet QikQikChannelRet
}

func LrNewQikQik(channelUserId, token string, args *map[string]interface{}) *QikQik {
	ret := new(QikQik)
	ret.Init(channelUserId, token, args)
	return ret
}

func (this *QikQik) Init(channelUserId, token string, args *map[string]interface{}) {
	this.Lr.Init(channelUserId, token)
	appid := (*args)["APPID"].(string)
	appsecret := (*args)["APP_SECRET"].(string)
	signContest := channelUserId + token  + appid + appsecret
	beego.Trace(signContest)
	sign := tool.Md5([]byte(signContest))
	format := "http://api.sy.7k7k.com/index.php/user/checkUser/uid/%s/vkey/%s/appid/%s/sign/%s"
	this.Url = fmt.Sprintf(format, channelUserId, token, appid, sign)
	beego.Trace(fmt.Sprintf(format, channelUserId, token, appid, sign))
}

func (this *QikQik) ParseChannelRet() (err error) {
	beego.Trace(this.Result)
	return json.Unmarshal([]byte(this.Result), &this.channelRet)
}

func (this *QikQik) CheckChannelRet() bool {
	return this.channelRet.Status == 0
}
