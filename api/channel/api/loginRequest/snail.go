package loginRequest

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"u9/models"
	"u9/tool"
)

//当乐

type SnailChannelRet struct {
	ErrorCode   string 	`json:"ErrorCode"`
	ErrorDesc   string  `json:"ErrorDesc"`
	Account 	string  `json:"Account"`
}

type Snail struct {
	Lr
	channelRet SnailChannelRet
}

func LrNewSnail(mlr *models.LoginRequest, args *map[string]interface{}) *Snail {
	ret := new(Snail)
	ret.Init(mlr, args)
	return ret
}

func (this *Snail) Init(mlr *models.LoginRequest, args *map[string]interface{}) {
	this.Lr.Init(mlr)
	appid := (*args)["SNAIL_APPID"].(string)
	appkey := (*args)["SNAIL_APPKEY"].(string)
	sign := tool.Md5([]byte(appid  +"4"+  mlr.ChannelUserid + this.mlr.Token +appkey))
	format := "http://api.app.snail.com/store/platform/sdk/ap?AppId=%s&Act=4&Uin=%s&SessionId=%s&Sign=%s"
	this.Url = fmt.Sprintf(format, appid, mlr.ChannelUserid,mlr.Token, sign)
 	beego.Trace(this.Url)
}

func (this *Snail) ParseChannelRet() (err error) {
	beego.Trace(this.Result)
	return json.Unmarshal([]byte(this.Result), &this.channelRet)
}

func (this *Snail) CheckChannelRet() bool {
	return this.channelRet.ErrorCode == "1"
}
