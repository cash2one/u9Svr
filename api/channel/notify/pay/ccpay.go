package channelPayNotify

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"net/url"
	"strconv"
	"u9/tool"
)

var ccpayUrlKeys []string = []string{"transactionNo", "partnerTransactionNo", "statusCode",
	"productId", "orderPrice", "packageId", "sign"}

const (
	err_ccpayParseAppSecret = 10401
	err_ccpayResultFailure  = 10402
)

type CCPay struct {
	Base
	appSecret string
}

func NewCCPay(channelId, productId int, urlParams *url.Values) *CCPay {
	ret := new(CCPay)
	ret.Init(channelId, productId, urlParams)
	return ret
}

func (this *CCPay) Init(channelId, productId int, urlParams *url.Values) {
	this.Base.Init(channelId, productId, urlParams, &ccpayUrlKeys)
}

func (this *CCPay) parseAppSecret() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_ccpayParseAppSecret
			beego.Trace(err)
		}
	}()
	this.appSecret, err = this.getPackageParam("app_secret")
	return
}

func (this *CCPay) parseUrlParam() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_parseUrlParam
			beego.Trace(err)
		}
	}()

	this.orderId = this.urlParams.Get("partnerTransactionNo")
	this.channelOrderId = this.urlParams.Get("transactionNo")

	payAmount := 0.0
	if payAmount, err = strconv.ParseFloat(this.urlParams.Get("orderPrice"), 64); err != nil {
		return err
	} else {
		this.payAmount = int(payAmount * 100)
	}

	return
}

func (this *CCPay) ParseParam() (err error) {
	if err = this.parseAppSecret(); err != nil {
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

func (this *CCPay) CheckSign() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_checkSign
			beego.Trace(err)
		}
	}()

	format := "orderPrice=%s&packageId=%s&partnerTransactionNo=%s&productId=%s&statusCode=%s&transactionNo=%s&%s"

	urlSign := this.urlParams.Get("sign")
	content := fmt.Sprintf(format, this.urlParams.Get("orderPrice"), this.urlParams.Get("packageId"),
		this.urlParams.Get("partnerTransactionNo"), this.urlParams.Get("productId"), this.urlParams.Get("statusCode"),
		this.urlParams.Get("transactionNo"), this.appSecret)

	if sign := tool.Md5([]byte(content)); sign != urlSign {
		msg := fmt.Sprintf("Sign is invalid, content:%s, sign:%s, urlSign%s", content, sign, urlSign)
		err = errors.New(msg)
		return
	}
	return
}

func (this *CCPay) ParseChannelRet() (err error) {
	if result := this.urlParams.Get("statusCode"); result != "0000" {
		this.callbackRet = err_ccpayResultFailure
	}
	return
}

func (this *CCPay) GetResult() (ret string) {
	if this.callbackRet == err_noerror {
		ret = "success"
	} else {
		beego.Trace(this.urlParams)
		ret = "fail"
	}
	return
}
