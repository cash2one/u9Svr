package channelPayNotify

import (
	"fmt"
)

var youlongUrlKeys []string = []string{"orderId", "userName", "amount", "flag"}

//游龙
type YouLong struct {
	MD5
}

func (this *YouLong) Init(params ...interface{}) (err error) {
	if err = this.MD5.Init(params...); err != nil {
		return
	}

	this.urlParamCheckKeys = &youlongUrlKeys
	this.requireChannelUserId = true
	this.exChangeRatio = 100

	this.channelParamKeys["_gameId"] = ""
	this.channelParamKeys["_gameKey"] = ""
	this.channelParamKeys["_payKey"] = "PKEY"

	this.channelTradeData = nil
	this.channelRetData = nil

	this.signHandleMethod = "ToUpper"
	return
}

func (this *YouLong) ParseInputParam(params ...interface{}) (err error) {
	if err = this.MD5.ParseInputParam(); err != nil {
		return
	}

	orderId_key := "orderId"
	channelUserId_key := "userName"
	channelOrderId_key := ""
	amount_key := "amount"
	discount_key := ""
	return parseTradeData_urlParam(&this.MD5.Base,
		orderId_key, channelUserId_key, channelOrderId_key,
		amount_key, discount_key)
}

func (this *YouLong) CheckSign(params ...interface{}) (err error) {
	format := "%s%s%s%s%s"
	this.signContent = fmt.Sprintf(format,
		this.orderId,
		this.channelUserId,
		this.urlParams.Get("amount"),
		this.urlParams.Get("extra"),
		this.channelParams["_payKey"])
	this.inputSign = this.urlParams.Get("flag")
	return this.MD5.CheckSign()
}

func (this *YouLong) GetResult(params ...interface{}) (ret string) {
	format := `%s`
	succMsg := "OK"
	failMsg := "NO"
	return this.MD5.GetResult(format, succMsg, failMsg)
}
