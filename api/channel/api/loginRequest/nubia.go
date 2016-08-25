package loginRequest

import (
	"fmt"
	"encoding/json"
	"github.com/astaxie/beego"
	"strconv"
	"time"
	"u9/models"
	"u9/tool"
)

//努比亚

type NubiaChannelRet struct {
	Code    int     `json:"code"`
	// data   	string 	`json:"data"`
	Message string  `json:"message"`
}

type Nubia struct {
	Lr
	channelRet NubiaChannelRet
}

func LrNewNubia(mlr *models.LoginRequest, args *map[string]interface{}) *Nubia {
	ret := new(Nubia)
	ret.Init(mlr, args)
	return ret
}

func (this *Nubia) Init(mlr *models.LoginRequest, args *map[string]interface{}) {
	this.Lr.Init(mlr)
	appid := (*args)["NUBIA_APPID"].(string)
	secretKey := (*args)["NUBIA_APPSECRET_KEY"].(string)
	this.Method = "POST"
	t := strconv.FormatInt(time.Now().Unix(), 10)
	context := "data_timestamp=" + t + "&session_id=" + this.mlr.Token + "&uid=" + this.mlr.ChannelUserid + 
	 ":" + appid + ":" + secretKey
	sign := tool.Md5([]byte(context))
	format := "http://niusdk.api.nubia.cn/VerifyAccount/CheckSession?uid=%s&session_id=%s&data_timestamp=%s&sign=%s"
	this.Url = fmt.Sprintf(format, this.mlr.ChannelUserid, this.mlr.Token, t, sign)
	beego.Trace("url:",this.Url)
}

func (this *Nubia) ParseChannelRet() (err error) {
	beego.Trace(this.Result)
	return json.Unmarshal([]byte(this.Result), &this.channelRet)
}

func (this *Nubia) CheckChannelRet() bool {
	return this.channelRet.Code == 0
}
