package lrBefore

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
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
	base
	authorizationCode string
	clientId          string
	clientSecret      string
}

func NewXMW() *XMW {
	ret := new(XMW)
	return ret
}

func (this *XMW) Init(param *Param) (err error) {
	this.base.Init(param)
	this.channelRet = new(xmwChannelRet)

	defer func() {
		if err != nil {
			format := "exec: err:%v"
			msg := fmt.Sprintf(format, err) + this.Dump()
			err = errors.New(msg)
			beego.Error(err)
		}
	}()

	if this.clientId, err = this.getChannelParam("XMWAPPID"); err != nil {
		return
	}

	if this.clientSecret, err = this.getChannelParam("XMWAPPSECRET"); err != nil {
		return
	}
	return
}

func (this *XMW) getAccessToken() (err error) {
	this.Method = "POST"
	format := "http://open.xmwan.com/v2/oauth2/access_token?client_id=%s&client_secret=%s&grant_type=%s&code=%s"
	this.Url = fmt.Sprintf(format, this.clientId, this.clientSecret, "authorization_code", this.param.Token)
	this.base.InitParam()

	if err = this.GetResponse(); err != nil {
		beego.Error(err)
		return
	}

	beego.Trace("getAccessToken:", this.Result)
	if err = json.Unmarshal([]byte(this.Result), &this.channelRet); err != nil {
		beego.Error(err)
		return
	}

	channelRet := this.channelRet.(*xmwChannelRet)
	if channelRet.Error != "" {
		err = errors.New(channelRet.Error)
		beego.Error(err)
		return
	}
	this.param.Token = channelRet.AccessToken
	return
}

func (this *XMW) getUserInfo() (err error) {
	channelRet := this.channelRet.(*xmwChannelRet)

	this.Method = "GET"
	this.Url = "http://open.xmwan.com/v2/users/me?access_token=" + channelRet.AccessToken
	this.base.InitParam()

	if err = this.GetResponse(); err != nil {
		beego.Error(err)
		return
	}

	beego.Trace("getUserInfo:", this.Result)
	if err = json.Unmarshal([]byte(this.Result), &this.channelRet); err != nil {
		beego.Error(err)
		return
	}

	if channelRet.Error != "" {
		err = errors.New(channelRet.Error)
		beego.Error(err)
		return
	}

	this.param.ChannelUserId = channelRet.XmwOpenId
	this.param.ChannelUserName = channelRet.Nickname
	return
}

func (this *XMW) Exec() (ret string, err error) {
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
