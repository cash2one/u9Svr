package channelPayNotify

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"net/url"
	"strconv"
	"strings"
	"u9/tool"
)

var yijieUrlKeys []string = []string{"app", "ct", "fee", "sdk", "pt", "ssid",
	"tcd", "uid", "sign"}

const (
	err_yijieParsePayKey   = 13201
	err_yijieResultFailure = 13202
)

//易接
type YiJie struct {
	Base
	yijie_appId     string
	yijie_payKey    string
	yijie_channelId string
}

func NewYiJie(channelId, productId int, urlParams *url.Values) *YiJie {
	ret := new(YiJie)
	ret.Init(channelId, productId, urlParams)
	return ret
}

func (this *YiJie) Init(channelId, productId int, urlParams *url.Values) {
	this.Base.Init(channelId, productId, urlParams, &yijieUrlKeys)
}

func (this *YiJie) parsePayKey() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_yijieParsePayKey
			beego.Trace(err)
		}
	}()
	this.yijie_appId, err = this.getPackageParam("com.snowfish.appid")
	this.yijie_payKey, err = this.getPackageParam("com.snowfish.appsecret")
	this.yijie_channelId, err = this.getPackageParam("com.snowfish.channelid")
	this.yijie_appId = replaceSDKParam(this.yijie_appId)
	this.yijie_channelId = replaceSDKParam(this.yijie_channelId)
	return
}

func (this *YiJie) parseUrlParam() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_parseUrlParam
			beego.Trace(err)
		}
	}()

	this.orderId = this.urlParams.Get("cbi")
	this.channelUserId = this.urlParams.Get("uid")
	this.channelOrderId = this.urlParams.Get("tcd")
	var money string = this.urlParams.Get("fee")
	if this.payAmount, err = strconv.Atoi(money); err != nil {
		beego.Trace(err)
	}

	return
}

func (this *YiJie) ParseChannelRet() (err error) {
	if result := this.urlParams.Get("st"); result != "1" {
		this.callbackRet = err_yijieResultFailure
	}
	return
}

func (this *YiJie) ParseParam() (err error) {
	if err = this.parseUrlParam(); err != nil {
		return
	}
	if err = this.parsePayKey(); err != nil {
		return
	}
	if err = this.Base.ParseParam(); err != nil {
		return
	}
	return
}

func (this *YiJie) CheckSign() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_checkSign
			beego.Trace(err)
		}
	}()

	format := "app=%s&cbi=%s&ct=%s&fee=%s&pt=%s&sdk=%s&ssid=%s&st=%s&tcd=%s&uid=%s&ver=%s%s"
	content := fmt.Sprintf(format,
		strings.ToLower(this.yijie_appId), this.orderId, this.urlParams.Get("ct"),
		this.urlParams.Get("fee"), this.urlParams.Get("pt"), strings.ToLower(this.yijie_channelId),
		this.urlParams.Get("ssid"), this.urlParams.Get("st"), this.channelOrderId,
		this.channelUserId, this.urlParams.Get("ver"), this.yijie_payKey)

	urlSign := this.urlParams.Get("sign")
	if sign := tool.Md5([]byte(content)); sign != urlSign {
		beego.Trace("url:", this.urlParams)
		msg := fmt.Sprintf("Sign is invalid, content:%s, sign:%s, urlSign:%s", content, sign, urlSign)
		err = errors.New(msg)
		return
	}

	return
}

func (this *YiJie) GetResult() (ret string) {
	if this.callbackRet == err_noerror {
		ret = "SUCCESS"
	} else {
		ret = "FAILURE"
	}
	return
}

func replaceSDKParam(s string) (ret string) {
	ret = strings.Replace(s, "{", "", -1)
	ret = strings.Replace(ret, "}", "", -1)
	ret = strings.Replace(ret, "-", "", -1)
	return
}
