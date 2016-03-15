package loginRequestHandle

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"u9/models"
)

type XMWChannelRet struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	XmwOpenId    string `json:"xmw_open_id"`
	Nickname     string `json:"nickname"`
	Avatar       string `json:"avatar"`
	Gender       int    `json:"gender"`
}

type XMW struct {
	LRH
	authorizationCode string
	clientId          string
	clientSecret      string
	channelRet        XMWChannelRet
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
	this.Url = "http://open.xmwan.com/v2/oauth2/access_token"
	this.LRH.InitParam()

	this.Req.Param("client_id", this.clientId)
	this.Req.Param("client_secret", this.clientSecret)
	this.Req.Param("grant_type", "authorization_code")
	this.Req.Param("code", this.param.Token)

	if err = this.Response(); err != nil {
		beego.Error(err)
		return
	}

	beego.Trace("getAccessToken:", this.Result)
	if err = json.Unmarshal([]byte(this.Result), &this.channelRet); err != nil {
		beego.Error(err)
		return
	}
	return
}

func (this *XMW) getUserInfo() (err error) {
	this.Method = "GET"
	this.Url = "http://open.xmwan.com/v2/users/me?access_token=" + this.channelRet.AccessToken
	this.LRH.InitParam()

	if err = this.Response(); err != nil {
		beego.Error(err)
		return
	}

	beego.Trace("getUserInfo:", this.Result)
	if err = json.Unmarshal([]byte(this.Result), &this.channelRet); err != nil {
		beego.Error(err)
		return
	}

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
