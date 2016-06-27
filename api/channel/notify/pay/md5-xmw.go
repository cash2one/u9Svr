package channelPayNotify

import (
	"fmt"
)

var xmwUrlKeys []string = []string{"serial", "amount", "status", "app_order_id",
	"app_user_id", "sign"}

type Xmw struct {
	MD5
}

func (this *Xmw) Init(params ...interface{}) (err error) {
	if err = this.MD5.Init(params...); err != nil {
		return
	}

	this.urlParamCheckKeys = &xmwUrlKeys
	this.requireChannelUserId = true
	this.exChangeRatio = 100

	this.channelParamKeys["_gameId"] = ""
	this.channelParamKeys["_gameKey"] = ""
	this.channelParamKeys["_payKey"] = "XMWAPPSECRET"

	this.channelTradeData = nil
	this.channelRetData = nil

	this.signHandleMethod = ""
	return
}

func (this *Xmw) ParseInputParam(params ...interface{}) (err error) {
	if err = this.MD5.ParseInputParam(); err != nil {
		return
	}

	orderId_key := "app_order_id"
	channelUserId_key := "app_user_id"
	channelOrderId_key := "serial"
	amount_key := "amount"
	discount_key := ""
	return parseTradeData_urlParam(&this.MD5.Base,
		orderId_key, channelUserId_key, channelOrderId_key,
		amount_key, discount_key)
}

func (this *Xmw) CheckSign(params ...interface{}) (err error) {
	format := "amount=%s&app_order_id=%s&app_user_id=%s&serial=%s&status=%s&client_secret=%s"
	this.signContent = fmt.Sprintf(format,
		this.urlParams.Get("amount"),
		this.orderId, this.channelUserId,
		this.channelOrderId,
		this.urlParams.Get("status"),
		this.channelParams["_payKey"])
	this.inputSign = this.urlParams.Get("sign")
	return this.MD5.CheckSign()
}

func (this *Xmw) CheckChannelRet(params ...interface{}) (err error) {
	tradeState := this.urlParams.Get("status") == "success"
	tradeFailDesc := `urlParam(status)!="success"`
	return this.MD5.CheckChannelRet(tradeState, tradeFailDesc)
}

func (this *Xmw) GetResult(params ...interface{}) (ret string) {
	format := `%s`
	succMsg := "success"
	failMsg := "fail"
	return this.MD5.GetResult(format, succMsg, failMsg)
}
