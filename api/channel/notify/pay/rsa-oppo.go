package channelPayNotify

import (
	"fmt"
)

var oppoUrlKeys []string = []string{"notifyId", "partnerOrder", "price", "count",
	"sign"}

const oppoRsaPublicKey = `MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCmreYIkPwVovKR8rLHWlFVw7YDfm9uQOJKL89Smt6ypXGVdrAKKl0wNYc3/jecAoPi2ylChfa2iRu5gunJyNmpWZzlCNRIau55fxGW0XEu553IiprOZcaw5OuYGlf60ga8QT6qToP0/dpiL/ZbmNUO9kUhosIjEu22uFgR+5cYyQIDAQAB`

type Oppo struct {
	Rsa
}

func (this *Oppo) Init(params ...interface{}) (err error) {
	if err = this.Rsa.Init(params...); err != nil {
		return
	}

	this.urlParamCheckKeys = &oppoUrlKeys
	this.requireChannelUserId = false
	this.exChangeRatio = 1

	this.channelParamKeys["_gameId"] = ""
	this.channelParamKeys["_gameKey"] = ""
	this.channelParamKeys["_payKey"] = ""

	this.channelParams["_payKey"] = oppoRsaPublicKey

	this.channelTradeData = nil
	this.channelRetData = nil

	this.signMode = 0
	return
}

func (this *Oppo) ParseInputParam(params ...interface{}) (err error) {
	if err = this.Rsa.ParseInputParam(); err != nil {
		return
	}

	orderId_key := "partnerOrder"
	channelUserId_key := ""
	channelOrderId_key := "notifyId"
	amount_key := "price"
	discount_key := ""
	return parseTradeData_urlParam(&this.Rsa.Base,
		orderId_key, channelUserId_key, channelOrderId_key,
		amount_key, discount_key)
}

func (this *Oppo) CheckSign(params ...interface{}) (err error) {
	format := "notifyId=%s&partnerOrder=%s&productName=%s&productDesc=%s&price=%d&count=%s&attach=%s"
	this.signContent = fmt.Sprintf(format,
		this.channelOrderId,
		this.orderId,
		this.urlParams.Get("productName"),
		this.urlParams.Get("productDesc"),
		this.payAmount,
		this.urlParams.Get("count"),
		this.urlParams.Get("attach"))

	this.inputSign = this.urlParams.Get("sign")
	return this.Rsa.CheckSign()
}

func (this *Oppo) GetResult(params ...interface{}) (ret string) {
	format := `result=%s&resultMsg=%s`
	succMsg := "OK," + errorDescList[this.lastError]
	failMsg := "FAIL," + errorDescList[this.lastError]
	return this.Rsa.GetResult(format, succMsg, failMsg)
}
