package channelPayNotify

import (
	"crypto/rsa"
	"fmt"
	"github.com/astaxie/beego"
	"net/url"
	"strconv"
	"u9/tool"
)

var oppoUrlKeys []string = []string{"notifyId", "partnerOrder", "price", "count", "sign"}

const (
	err_oppoInitRsaPublicKey = 11001
)

const oppoRsaPublicKeyStr = `MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCmreYIkPwVovKR8rLHWlFVw7YDfm9uQOJKL89Smt6ypXGVdrAKKl0wNYc3/jecAoPi2ylChfa2iRu5gunJyNmpWZzlCNRIau55fxGW0XEu553IiprOZcaw5OuYGlf60ga8QT6qToP0/dpiL/ZbmNUO9kUhosIjEu22uFgR+5cYyQIDAQAB`

var (
	oppoRsaPublicKey *rsa.PublicKey
)

type Oppo struct {
	Base
}

func NewOppo(channelId, productId int, urlParams *url.Values) *Oppo {
	ret := new(Oppo)
	ret.Init(channelId, productId, urlParams)
	return ret
}

func (this *Oppo) Init(channelId, productId int, urlParams *url.Values) {
	this.Base.Init(channelId, productId, urlParams, &oppoUrlKeys)
	this.existChannelUserId = false
	this.initRsaPublicKey()
}

func (this *Oppo) initRsaPublicKey() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_oppoInitRsaPublicKey
			beego.Trace(err)
		}
	}()

	if oppoRsaPublicKey == nil {
		oppoRsaPublicKey, err = tool.ParsePKIXPublicKeyWithStr(oppoRsaPublicKeyStr)
		if err != nil {
			beego.Error(err)
			return err
		}
	}
	return nil
}

func (this *Oppo) ParseParam() (err error) {
	if err = this.parseUrlParam(); err != nil {
		return
	}
	if err = this.Base.ParseParam(); err != nil {
		return
	}
	return
}

func (this *Oppo) parseUrlParam() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_parseUrlParam
			beego.Trace(err)
		}
	}()

	this.orderId = this.urlParams.Get("partnerOrder")
	this.channelOrderId = this.urlParams.Get("notifyId")

	payAmount := 0
	if payAmount, err = strconv.Atoi(this.urlParams.Get("price")); err != nil {
		return
	} else {
		this.payAmount = payAmount
	}
	return
}

func (this *Oppo) CheckSign() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_checkSign
			beego.Trace(err)
		}
	}()

	format := "notifyId=%s&partnerOrder=%s&productName=%s&productDesc=%s&price=%d&count=%s&attach=%s"
	context := fmt.Sprintf(format,
		this.channelOrderId, this.orderId, this.urlParams.Get("productName"),
		this.urlParams.Get("productDesc"), this.payAmount, this.urlParams.Get("count"),
		this.urlParams.Get("attach"))

	sign := this.urlParams.Get("sign")
	if err = tool.RsaVerifyPKCS1v15(oppoRsaPublicKey, context, sign); err != nil {
		msg := fmt.Sprintf("RsaVerifyPK CS1v15 exception: context:%s, sign:%s", context, sign)
		beego.Trace(msg)
		return err
	}
	return nil
}

func (this *Oppo) GetResult() (ret string) {
	result := ""
	if this.callbackRet == err_noerror {
		result = "OK"
	} else {
		result = "FAIL"
	}
	return fmt.Sprintf("result=%s&resultMsg=%s", result, strconv.Itoa(this.callbackRet))
}
