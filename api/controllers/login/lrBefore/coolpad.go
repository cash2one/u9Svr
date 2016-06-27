package lrBefore

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
)

type coolpadChannelRet struct {
	AccessToken       string `json:"access_token"`
	ExpiresIn         string `json:"expires_in`
	RefreshToken      string `json:"refresh_token`
	Openid            string `json:"openid"`
	Error             string `json:"error"`
	Error_description string `json:"error_description"`
}

type CoolPad struct {
	base
}

func NewCoolPad() *CoolPad {
	ret := new(CoolPad)
	return ret
}

func (this *CoolPad) Init(param *Param) (err error) {
	this.base.Init(param)
	this.channelRet = new(coolpadChannelRet)
	return
}

func (this *CoolPad) Exec() (ret string, err error) {
	this.Method = "GET"
	this.IsHttps = false

	defer func() {
		if err != nil {
			format := "exec: err:%v"
			msg := fmt.Sprintf(format, err) + this.Dump()
			err = errors.New(msg)
			beego.Error(err)
		}
	}()

	clientId := ""
	if clientId, err = this.getChannelParam("COOLPAD_APPID"); err != nil {
		return
	}

	appKey := ""
	if appKey, err = this.getChannelParam("COOLPAD_APPKEY"); err != nil {
		return
	}

	format := `https://openapi.coolyun.com/oauth2/token?` +
		`grant_type=authorization_code&client_id=%s&client_secret=%s` +
		`&code=%s&redirect_uri=%s`
	this.Url = fmt.Sprintf(format, clientId, appKey, this.param.Token, appKey)

	this.base.InitParam()
	if err = this.GetResponse(); err != nil {
		beego.Error(err)
		return
	}

	if err = json.Unmarshal([]byte(this.Result), &this.channelRet); err != nil {
		beego.Error(err)
		return
	}

	channelRet := this.channelRet.(*coolpadChannelRet)

	if channelRet.Error != "" {
		msg := fmt.Sprintf(`channelRet.Error!=""`)
		err = errors.New(msg)
		return
	}

	this.param.ChannelUserId = channelRet.Openid
	//this.param.ChannelUserName = channelRet.Openid

	data, _ := json.Marshal(this.channelRet)

	ret = string(data)
	return
}
