package channelPayNotify

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"net/url"
	"strconv"
	"u9/tool"
)

var haimawanUrlKeys []string = []string{"notify_time", "appid", "out_trade_no", "total_fee",
	"subject", "body", "trade_status", "sign"}

const (
	err_haimawanParseServerKey = 12501
	err_haimawanResultFailure  = 12502
)

//海马玩
type Haimawan struct {
	Base
	serverKey string
}

func NewHaimawan(channelId, productId int, urlParams *url.Values) *Haimawan {
	ret := new(Haimawan)
	ret.Init(channelId, productId, urlParams)
	return ret
}

func (this *Haimawan) Init(channelId, productId int, urlParams *url.Values) {
	this.Base.Init(channelId, productId, urlParams, &haimawanUrlKeys)
}

func (this *Haimawan) parseServerKey() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_haimawanParseServerKey
			beego.Trace(err)
		}
	}()
	this.serverKey, err = this.getPackageParam("HM_SERVERKEY")
	return
}

func (this *Haimawan) parseUrlParam() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_parseUrlParam
			beego.Trace(err)
		}
	}()

	this.orderId = this.urlParams.Get("out_trade_no")

	payAmount := 0.0
	if payAmount, err = strconv.ParseFloat(this.urlParams.Get("total_fee"), 64); err != nil {
		return err
	} else {
		this.payAmount = int(payAmount * 100)
	}
	return
}

func (this *Haimawan) ParseChannelRet() (err error) {
	if result := this.urlParams.Get("trade_status"); result != "1" {
		this.callbackRet = err_haimawanResultFailure
	}
	return
}

func (this *Haimawan) ParseParam() (err error) {
	if err = this.parseUrlParam(); err != nil {
		return
	}
	if err = this.parseServerKey(); err != nil {
		return
	}
	if err = this.Base.ParseParam(); err != nil {
		return
	}
	this.channelUserId = this.loginRequest.ChannelUserid
	this.channelOrderId = this.orderRequest.ChannelOrderId
	return
}

func (this *Haimawan) CheckSign() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_checkSign
			beego.Trace(err)
		}
	}()

	format := "notify_time=%s&appid=%s&out_trade_no=%s&total_fee=%s&subject=%s&body=%s&trade_status=%s%s"

	context := fmt.Sprintf(format,
		url.QueryEscape(this.urlParams.Get("notify_time")),
		url.QueryEscape(this.urlParams.Get("appid")),
		url.QueryEscape(this.urlParams.Get("out_trade_no")),
		url.QueryEscape(this.urlParams.Get("total_fee")),
		url.QueryEscape(this.urlParams.Get("subject")),
		url.QueryEscape(this.urlParams.Get("body")),
		url.QueryEscape(this.urlParams.Get("trade_status")), this.serverKey)

	if sign := tool.Md5([]byte(context)); sign != this.urlParams.Get("sign") {
		msg := fmt.Sprintf("Sign is invalid, context:%s, sign:%s", context, sign)
		err = errors.New(msg)
		return
	}
	return
}

func (this *Haimawan) GetResult() (ret string) {
	if this.callbackRet == err_noerror {
		ret = "success"
	} else {
		ret = "failure"
	}
	return
}
