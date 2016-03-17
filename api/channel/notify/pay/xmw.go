package channelPayNotify

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"net/url"
	"strconv"
	"u9/tool"
)

var xmwUrlKeys []string = []string{"serial", "amount", "status", "app_order_id",
	"app_user_id", "sign"}

const (
	err_xmwParseAppSecret = 12301
	err_xmwResultFailure  = 12302
)

type Xmw struct {
	Base
	appSecret string
}

func NewXmw(channelId, productId int, urlParams *url.Values) *Xmw {
	ret := new(Xmw)
	ret.Init(channelId, productId, urlParams)
	return ret
}

func (this *Xmw) Init(channelId, productId int, urlParams *url.Values) {
	this.Base.Init(channelId, productId, urlParams, &xmwUrlKeys)
}

func (this *Xmw) parseAppSecret() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_xmwParseAppSecret
			beego.Trace(err)
		}
	}()
	this.appSecret, err = this.getPackageParam("XMWAPPSECRET")
	return
}

func (this *Xmw) parseUrlParam() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_parseUrlParam
			beego.Trace(err)
		}
	}()

	this.orderId = this.urlParams.Get("app_order_id")
	this.channelUserId = this.urlParams.Get("app_user_id")
	this.channelOrderId = this.urlParams.Get("serial")

	payAmount := 0.0
	if payAmount, err = strconv.ParseFloat(this.urlParams.Get("amount"), 64); err != nil {
		return err
	} else {
		this.payAmount = int(payAmount * 100)
	}
	return
}

func (this *Xmw) ParseChannelRet() (err error) {
	if result := this.urlParams.Get("status"); result != "success" {
		this.callbackRet = err_xmwResultFailure
	}
	return
}

func (this *Xmw) ParseParam() (err error) {
	if err = this.parseUrlParam(); err != nil {
		return
	}
	if err = this.parseAppSecret(); err != nil {
		return
	}
	if err = this.Base.ParseParam(); err != nil {
		return
	}
	return
}

func (this *Xmw) CheckSign() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_checkSign
			beego.Trace(err)
		}
	}()

	format := "amount=%s&app_order_id=%s&app_user_id=%s&serial=%s&status=%s&client_secret=%s"
	context := fmt.Sprintf(format,
		this.urlParams.Get("amount"),
		this.orderId, this.channelUserId,
		this.channelOrderId, this.urlParams.Get("status"), this.appSecret)

	if sign := tool.Md5([]byte(context)); sign != this.urlParams.Get("sign") {
		msg := fmt.Sprintf("Sign is invalid, context:%s, sign:%s", context, sign)
		err = errors.New(msg)
		return
	}
	return
}

func (this *Xmw) GetResult() (ret string) {
	if this.callbackRet == err_noerror {
		ret = "success"
	} else {
		ret = "fail"
	}
	return
}
