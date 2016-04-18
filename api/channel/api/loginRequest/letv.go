package loginRequest

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"u9/models"
	// "u9/tool"
)

//乐视

type LeTVChannelRet struct {
	Request string     `json:"request"`
	Result  LeTvResult `json:"result"`
	Status  int        `json:"status"`
}
type LeTvResult struct {
	Letv_uid string `json:letv_uid`
	Nickname string `json:nickname`
	File300  string `json:file_300*300`
	File200  string `json:file_200*200`
	File70   string `json:file_70*70`
	File50   string `json:file_50*50`
}

type LeTV struct {
	Lr
	channelRet LeTVChannelRet
}

func LrNewLeTV(mlr *models.LoginRequest, args *map[string]interface{}) *LeTV {
	ret := new(LeTV)
	ret.Init(mlr, args)
	return ret
}

func (this *LeTV) Init(mlr *models.LoginRequest, args *map[string]interface{}) {
	this.Lr.Init(mlr)
	appid := (*args)["lepay_appid"].(string)
	format := "https://sso.letv.com/oauthopen/userbasic?" +
		"client_id=%s&uid=%s&access_token=%s"
	this.Url = fmt.Sprintf(format, appid, this.mlr.ChannelUserid, this.mlr.Token)
	beego.Trace(this.Url)
}

func (this *LeTV) ParseChannelRet() (err error) {
	beego.Trace(this.Result)
	return json.Unmarshal([]byte(this.Result), &this.channelRet)
}

func (this *LeTV) CheckChannelRet() bool {
	return this.channelRet.Status == 1
}
