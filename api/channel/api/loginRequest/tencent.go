package loginRequest

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"strconv"
	"time"
	"u9/models"
	"u9/tool"
)

type TencentLrChannelRet struct {
	Ret int    `json:"ret"`
	Msg string `json:"msg"`
}

type TencentExtParam struct {
	LoginType string `json:"loginType"`
	Debug     bool   `json:"debug"`
}

type Tencent struct {
	Lr
	channelRet TencentLrChannelRet
	args       *map[string]interface{}
	appId      string
	appKey     string
}

func LrNewTencent(mlr *models.LoginRequest, args *map[string]interface{}) *Tencent {
	ret := new(Tencent)
	ret.Init(mlr, args)
	return ret
}

func (this *Tencent) Init(mlr *models.LoginRequest, args *map[string]interface{}) {
	this.Lr.Init(mlr)
	this.args = args

	var extParam TencentExtParam
	json.Unmarshal([]byte(this.mlr.Ext), &extParam)

	if extParam.LoginType == "QQ" {
		if extParam.Debug {
			this.Url = "http://ysdk.qq.com/auth/qq_check_token"
		} else {
			this.Url = "http://ysdktest.qq.com/auth/qq_check_token"
		}
		this.appId = (*this.args)["QQ_APP_ID"].(string)
		this.appKey = (*this.args)["QQ_APP_KEY"].(string)
	} else if extParam.LoginType == "WEIXIN" {
		if extParam.Debug {
			this.Url = "http://ysdktest.qq.com/auth/wx_check_token"
		} else {
			this.Url = "ysdk.qq.com/auth/wx_check_token"
		}
		this.appId = (*this.args)["WX_APP_ID"].(string)
		this.appKey = (*this.args)["WX_APP_KEY"].(string)
	} else {
		beego.Error(errors.New("login type is error, must in (QQ, WEIXIN)"))
	}
}

func (this *Tencent) InitParam() {
	this.Lr.InitParam()

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	sig := tool.Md5([]byte(this.appKey + timestamp))
	format := "?timestamp=%s&appid=%s&sig=%s&openid=%s&openkey=%s"
	context := fmt.Sprintf(format, timestamp, this.appId, sig, this.mlr.ChannelUserid, this.mlr.Token)
	this.Url = this.Url + context
}

func (this *Tencent) ParseChannelRet() (err error) {
	beego.Trace(this.Result)
	return json.Unmarshal([]byte(this.Result), &this.channelRet)
}

func (this *Tencent) CheckChannelRet() bool {
	return this.channelRet.Ret == 0
}
