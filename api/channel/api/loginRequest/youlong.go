package loginRequest

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
)

//游龙

type YouLongChannelRet struct {
	State    string `json:"state"`
	Username string `json:"username"`
	ErrMsg   string `json:"errMsg"`
	ErrcMsg  string `json:"errcMsg"`
}

type YouLong struct {
	Lr
	channelRet YouLongChannelRet
}

func LrNewYouLong(channelUserId, token string, args *map[string]interface{}) *YouLong {
	ret := new(YouLong)
	ret.Init(channelUserId, token, args)
	return ret
}

func (this *YouLong) Init(channelUserId, token string, args *map[string]interface{}) {
	this.Lr.Init(channelUserId, token)
	this.Method = "POST"
	appid := (*args)["PID"].(string)
	format := "http://ucapi.411game.com/Api/checkToken?token=%s&pid=%s&ip=%s"
	this.Url = fmt.Sprintf(format, token, appid, channelUserId)
}

func (this *YouLong) ParseChannelRet() (err error) {
	beego.Trace(this.Result)
	return json.Unmarshal([]byte(this.Result), &this.channelRet)
}

func (this *YouLong) CheckChannelRet() bool {
	return this.channelRet.State == "1"
}
