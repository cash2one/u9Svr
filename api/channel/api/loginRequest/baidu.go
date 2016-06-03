package loginRequest

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"net/url"
	"strconv"
	"u9/models"
	"u9/tool"
)

type baiduChannelRet struct {
	ResultCode int    `json:"ResultCode"`
	ResultMsg  string `json:"ResultMsg"`
	AppID      string `json:"AppID"`
	Sign       string `json:"Sign"`
	Content    string `json:"Content"`
}

type Baidu struct {
	Lr
	appId      string
	secretKey  string
	channelRet baiduChannelRet
}

func LrNewBaidu(mlr *models.LoginRequest, args *map[string]interface{}) *Baidu {
	ret := new(Baidu)
	ret.Init(mlr, args)
	return ret
}

func (this *Baidu) Init(mlr *models.LoginRequest, args *map[string]interface{}) {
	this.Lr.Init(mlr)
	this.appId = (*args)["BAIDU_APPID"].(string)
	this.secretKey = (*args)["BAIDU_APPSECRET"].(string)
	sign := tool.Md5([]byte(this.appId + this.mlr.Token + this.secretKey))

	this.Url = "http://querysdkapi.baidu.com/query/cploginstatequery?"
	format := "AppID=%s&AccessToken=%s&Sign=%s"
	this.Url = this.Url + fmt.Sprintf(format, this.appId, this.mlr.Token, sign)
	beego.Trace(this.Url)
}

func (this *Baidu) ParseChannelRet() (err error) {
	beego.Trace(this.Result)
	if err = json.Unmarshal([]byte(this.Result), &this.channelRet); err != nil {
		return
	}

	if this.appId != this.channelRet.AppID {
		err = errors.New("channel.AppId isn't equal with xmlParms.AppId.")
		return
	}

	content := this.appId + strconv.Itoa(this.channelRet.ResultCode) + this.channelRet.Content + this.secretKey
	if sign := tool.Md5([]byte(content)); sign != this.channelRet.Sign {
		msg := fmt.Sprintf("sign(%s) is equal channelRet's sign(%s)", sign, this.channelRet.Sign)
		err = errors.New(msg)
		return
	}

	if this.channelRet.ResultMsg, err = url.QueryUnescape(this.channelRet.ResultMsg); err != nil {
		return
	}

	var enByte []byte
	if enByte, err = base64.StdEncoding.DecodeString(this.channelRet.Content); err != nil {
		return
	}
	this.channelRet.Content = string(enByte)

	//beego.Trace(fmt.Sprintf("%+v",this.channelRet))
	return
}

func (this *Baidu) CheckChannelRet() bool {
	return this.channelRet.ResultCode == 1
}
