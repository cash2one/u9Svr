package channelPayNotify

import (
	"crypto/rsa"
	"fmt"
	"github.com/astaxie/beego"
	"net/url"
	"strconv"
	"u9/tool"
)

var amigoUrlKeys []string = []string{"api_key", "close_time", "create_time", "deal_price", "out_order_no",
	"pay_channel", "submit_time", "user_id", "sign"}

const (
	err_amigoParseApiKey      = 12001
	err_amigoResultFailure    = 12002
	err_amigoInitRsaPublicKey = 12003
	amigoRsaPublicKeyStr      = `MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCkdZLJ5OCrpxM9EdfTUPA/kM95dgcf7pNMvLmdSbXO9U/LVlhNg1q1EBABXzddK5kURM3vNShsfuAichOVJj+0rV3iYcdym9ZJA6cbRhwBWY76PMmfl9ysj+2g7DxIpNrA7mx0XEEC5++67meSO77qafnSRa884BHEBJF/RoGSBwIDAQAB`
)

var (
	amigoRsaPublicKey *rsa.PublicKey
)

type Amigo struct {
	Base
	apiKey string
}

func NewAmigo(channelId, productId int, urlParams *url.Values) *Amigo {
	ret := new(Amigo)
	ret.Init(channelId, productId, urlParams)
	return ret
}

func (this *Amigo) Init(channelId, productId int, urlParams *url.Values) {
	this.Base.Init(channelId, productId, urlParams, &amigoUrlKeys)
}

func (this *Amigo) initRsaPublicKey() (err error) {
	if amigoRsaPublicKey == nil {
		amigoRsaPublicKey, err = tool.ParsePKIXPublicKeyWithStr(amigoRsaPublicKeyStr)
		if err != nil {
			this.callbackRet = err_amigoInitRsaPublicKey
			beego.Error(err)
			return err
		}
	}
	return nil
}

func (this *Amigo) parseApiKey() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_amigoParseApiKey
			beego.Trace(err)
		}
	}()
	this.apiKey, err = this.getPackageParam("AMIGO_APIKEY")
	return
}

func (this *Amigo) parseUrlParam() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_parseUrlParam
			beego.Trace(err)
			beego.Trace(this.urlParams)
		}
	}()

	this.orderId = this.urlParams.Get("out_order_no")

	payAmount := 0.0
	if payAmount, err = strconv.ParseFloat(this.urlParams.Get("deal_price"), 64); err != nil {
		return err
	} else {
		this.payAmount = int(payAmount * 100)
	}
	return
}

func (this *Amigo) ParseChannelRet() (err error) {
	if this.orderId != this.orderRequest.OrderId {
		this.callbackRet = err_orderIsNotExist
		return
	}

	if this.orderRequest.ReqAmount != this.payAmount {
		this.callbackRet = err_payAmountError
	}

	return
}

func (this *Amigo) ParseParam() (err error) {
	if err = this.initRsaPublicKey(); err != nil {
		return
	}
	if err = this.parseUrlParam(); err != nil {
		return
	}
	if err = this.parseApiKey(); err != nil {
		return
	}
	if err = this.Base.ParseParam(); err != nil {
		return
	}
	this.channelUserId = this.loginRequest.ChannelUserid
	this.channelOrderId = this.orderRequest.ChannelOrderId
	return
}

func (this *Amigo) CheckSign() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_checkSign
			beego.Trace(err)
		}
	}()

	format := "api_key=%s&close_time=%s&create_time=%s&deal_price=%s&out_order_no=%s&pay_channel=%s&submit_time=%s&user_id=%s"
	context := fmt.Sprintf(format,
		this.apiKey, this.urlParams.Get("close_time"),
		this.urlParams.Get("create_time"), this.urlParams.Get("deal_price"),
		this.orderId, this.urlParams.Get("pay_channel"),
		this.urlParams.Get("submit_time"), "null")
	sign := this.urlParams.Get("sign")

	if err = tool.RsaVerifyPKCS1v15(amigoRsaPublicKey, context, sign); err != nil {
		msg := fmt.Sprintf("RsaVerifyPK CS1v15 exception: context:%s, sign:%s", context, sign)
		beego.Trace(msg)
		return err
	}
	return
}

func (this *Amigo) GetResult() (ret string) {
	beego.Trace(this.callbackRet)
	if this.callbackRet == err_noerror {
		ret = "success"
	} else {
		ret = "failure"
	}
	return
}
