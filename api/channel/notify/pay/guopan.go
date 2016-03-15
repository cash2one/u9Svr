package channelPayNotify

import (
	"github.com/astaxie/beego"
	"net/url"
	"strconv"
	"u9/tool"
)

var guopanUrlKeys []string = []string{"trade_no", "serialNumber",
	"money", "status", "t", "sign", "reserved"}

const (
	err_guopanParseSecrectKey = 10701
	err_guopanCheckSign       = 10702
	err_guopanCallbackFail    = 10703
)

type Guopan struct {
	Base
	secrectKey string
}

func NewGuopan(channelId, productId int, urlParams *url.Values) *Guopan {
	ret := new(Guopan)
	ret.Init(channelId, productId, urlParams)
	return ret
}

func (this *Guopan) Init(channelId, productId int, urlParams *url.Values) {
	this.Base.Init(channelId, productId, urlParams, &guopanUrlKeys)
}

func (this *Guopan) parseSecrectKey() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_guopanParseSecrectKey
			beego.Trace(err)
		}
	}()
	this.secrectKey, err = this.getPackageParam("GUOPAN_SERVER_SECRETKEY")
	return
}

func (this *Guopan) parseUrlParam() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_parseUrlParam
			beego.Trace(err)
		}
	}()

	this.orderId = this.urlParams.Get("serialNumber")
	this.channelOrderId = this.urlParams.Get("trade_no")
	// this.channelUserId = this.
	payAmount := 0.0
	if payAmount, err = strconv.ParseFloat(this.urlParams.Get("money"), 64); err != nil {
		return err
	} else {
		this.payAmount = int(payAmount * 100)
	}
	return
}

func (this *Guopan) ParseChannelRet() (err error) {
	if this.orderId != this.orderRequest.OrderId {
		this.callbackRet = err_orderIsNotExist
		return
	}

	if this.orderRequest.ReqAmount != this.payAmount {
		this.callbackRet = err_payAmountError
	}

	return
}

func (this *Guopan) ParseParam() (err error) {
	if err = this.parseSecrectKey(); err != nil {
		return
	}
	if err = this.parseUrlParam(); err != nil {
		return
	}
	if err = this.Base.ParseParam(); err != nil {
		return
	}
	this.channelUserId = this.loginRequest.ChannelUserid
	return
}

func (this *Guopan) CheckSign() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_checkSign
			beego.Trace(err)
		}
	}()
	serialNumber := this.urlParams.Get("serialNumber")
	status := this.urlParams.Get("status")
	if status != "1" {
		this.callbackRet = err_guopanCallbackFail
	}
	t := this.urlParams.Get("t")
	money := this.urlParams.Get("money")
	content := serialNumber + money + status + t + this.secrectKey
	sign := tool.Md5([]byte(content))
	beego.Trace(content)
	beego.Trace(sign)

	if sign != this.urlParams.Get("sign") {
		this.callbackRet = err_guopanCheckSign
	}
	return
}

func (this *Guopan) GetResult() (ret string) {
	if this.callbackRet == err_noerror {
		ret = "success"
	} else {
		ret = "fail"
	}
	return
}
