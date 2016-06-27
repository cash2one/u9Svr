package channelPayNotify

import (
	"fmt"
)

var yijieUrlKeys []string = []string{"app", "cbi", "ct", "fee", "pt",
	"sdk", "pt", "ssid", "st", "tcd", "uid", "ver", "sign"}

//易接
type YiJie struct {
	MD5
}

func (this *YiJie) Init(params ...interface{}) (err error) {
	if err = this.MD5.Init(params...); err != nil {
		return
	}

	this.urlParamCheckKeys = &yijieUrlKeys
	this.requireChannelUserId = true
	this.exChangeRatio = 1

	this.channelParamKeys["_gameId"] = ""
	this.channelParamKeys["_gameKey"] = ""
	this.channelParamKeys["_payKey"] = "com.snowfish.appsecret"

	this.channelTradeData = nil
	this.channelRetData = nil

	this.signHandleMethod = ""
	return
}

func (this *YiJie) ParseInputParam(params ...interface{}) (err error) {
	if err = this.MD5.ParseInputParam(); err != nil {
		return
	}

	orderId_key := "cbi"
	channelUserId_key := "uid"
	channelOrderId_key := "tcd"
	amount_key := "fee"
	discount_key := ""
	return parseTradeData_urlParam(&this.Base,
		orderId_key, channelUserId_key, channelOrderId_key,
		amount_key, discount_key)
}

func (this *YiJie) CheckSign(params ...interface{}) (err error) {
	format := "app=%s&cbi=%s&ct=%s&fee=%s&pt=%s&sdk=%s&ssid=%s&st=%s&tcd=%s" +
		"&uid=%s&ver=%s%s"
	this.signContent = fmt.Sprintf(format,
		this.urlParams.Get("app"),
		this.orderId,
		this.urlParams.Get("ct"),
		this.urlParams.Get("fee"),
		this.urlParams.Get("pt"),
		this.urlParams.Get("sdk"),
		this.urlParams.Get("ssid"),
		this.urlParams.Get("st"),
		this.channelOrderId,
		this.channelUserId,
		this.urlParams.Get("ver"),
		this.channelParams["_payKey"])
	this.inputSign = this.urlParams.Get("sign")
	return this.MD5.CheckSign()
}

func (this *YiJie) CheckChannelRet(params ...interface{}) (err error) {
	tradeState := this.urlParams.Get("st") == "1"
	tradeFailDesc := `urlParam(st)!="1"`
	return this.MD5.CheckChannelRet(tradeState, tradeFailDesc)
}

func (this *YiJie) GetResult(params ...interface{}) (ret string) {
	format := `%s`
	succMsg := "SUCCESS"
	failMsg := "FAILURE"
	return this.MD5.GetResult(format, succMsg, failMsg)
}
