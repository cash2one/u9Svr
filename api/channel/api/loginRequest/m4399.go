package loginRequest

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"u9/models"
)

type m4399ChannelRet struct {
	Code    string     `json:"code"`
	Result  ResultJson `json:"result"`
	Message string     `json:"message"`
}

type ResultJson struct {
	Uid string `json:"uid"`
}

type M4399 struct {
	Lr
	channelRet m4399ChannelRet
}

func LrNewM4399(mlr *models.LoginRequest, args *map[string]interface{}) *M4399 {
	ret := new(M4399)
	ret.Init(mlr, args)
	return ret
}

func (this *M4399) Init(mlr *models.LoginRequest, args *map[string]interface{}) {
	this.Lr.Init(mlr)
	format := "http://m.4399api.com/openapi/oauth-check.html?state=%s&uid=%s"
	this.Url = fmt.Sprintf(format, this.mlr.Token, this.mlr.ChannelUserid)
}

func (this *M4399) ParseChannelRet() (err error) {
	beego.Trace(this.Result)
	err = json.Unmarshal([]byte(this.Result), &this.channelRet)
	return
}

func (this *M4399) CheckChannelRet() bool {
	return this.channelRet.Code == "100"
}
