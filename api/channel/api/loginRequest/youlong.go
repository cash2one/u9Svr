package loginRequest

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"u9/models"
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

func LrNewYouLong(mlr *models.LoginRequest, args *map[string]interface{}) *YouLong {
	ret := new(YouLong)
	ret.Init(mlr, args)
	return ret
}

func (this *YouLong) Init(mlr *models.LoginRequest, args *map[string]interface{}) {
	this.Lr.Init(mlr)
	this.Method = "POST"
	appid := (*args)["PID"].(string)
	format := "http://ucapi.411game.com/Api/checkToken?token=%s&pid=%s&ip=%s"
	this.Url = fmt.Sprintf(format, this.mlr.Token, appid, this.mlr.ChannelUserid)
}

func (this *YouLong) ParseChannelRet() (err error) {
	beego.Trace(this.Result)
	return json.Unmarshal([]byte(this.Result), &this.channelRet)
}

func (this *YouLong) CheckChannelRet() bool {
	return this.channelRet.State == "1"
}
