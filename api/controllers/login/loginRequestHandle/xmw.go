package loginRequestHandle

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"u9/models"
)

type xmwChannelRet struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
	AccessToken      string `json:"access_token"`
	RefreshToken     string `json:"refresh_token"`
	ExpiresIn        int    `json:"expires_in"`
	XmwOpenId        string `json:"xmw_open_id"`
	Nickname         string `json:"nickname"`
	Avatar           string `json:"avatar"`
	Gender           int    `json:"gender"`
}

type XMW struct {
	LRH
	channelRet        xmwChannelRet
	authorizationCode string
	clientId          string
	clientSecret      string
}

func NewXMW() *XMW {
	ret := new(XMW)
	return ret
}

func (this *XMW) Init(param *Param) (err error) {
	this.LRH.Init(param)
	pp := new(models.PackageParam)

	this.clientId, err = pp.GetXmlParam(this.param.ChannelId, this.param.ProductId, "XMWAPPID")
	if err != nil {
		beego.Error(err)
		return
	}

	this.clientSecret, err = pp.GetXmlParam(this.param.ChannelId, this.param.ProductId, "XMWAPPSECRET")
	if err != nil {
		beego.Error(err)
		return
	}
	return
}

func (this *XMW) getAccessToken() (err error) {
	this.Method = "POST"
	format := "http://open.xmwan.com/v2/oauth2/access_token?client_id=%s&client_secret=%s&grant_type=%s&code=%s"
	this.Url = fmt.Sprintf(format, this.clientId, this.clientSecret, "authorization_code", this.param.Token)
	this.LRH.InitParam()

	if err = this.GetResponse(); err != nil {
		beego.Error(err)
		return
	}

	beego.Trace("getAccessToken:", this.Result)
	if err = json.Unmarshal([]byte(this.Result), &this.channelRet); err != nil {
		beego.Error(err)
		return
	}

	if this.channelRet.Error != "" {
		err = errors.New(this.channelRet.Error)
		beego.Error(err)
		return
	}
	//this.channelRet.Error = ""
	this.param.Token = this.channelRet.AccessToken
	return
}

func (this *XMW) getUserInfo() (err error) {
	this.Method = "GET"
	this.Url = "http://open.xmwan.com/v2/users/me?access_token=" + this.channelRet.AccessToken
	this.LRH.InitParam()

	if err = this.GetResponse(); err != nil {
		beego.Error(err)
		return
	}

	beego.Trace("getUserInfo:", this.Result)
	if err = json.Unmarshal([]byte(this.Result), &this.channelRet); err != nil {
		beego.Error(err)
		return
	}

	if this.channelRet.Error != "" {
		err = errors.New(this.channelRet.Error)
		beego.Error(err)
		return
	}
	//this.channelRet.Error = ""

	this.param.ChannelUserId = this.channelRet.XmwOpenId
	this.param.ChannelUserName = this.channelRet.Nickname
	return
}

func (this *XMW) Handle() (ret string, err error) {
	if err = this.getAccessToken(); err != nil {
		return
	}
	if err = this.getUserInfo(); err != nil {
		return
	}

	data, _ := json.Marshal(this.channelRet)
	ret = string(data)
	return
}

// func (this *XMW) GetChannelResult() (ret interface{}) {
// 	//return this.channelRet.AccessToken
// 	return this.channelRet
// }
