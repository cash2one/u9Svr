package channelPayNotify

import (
	"fmt"
)

var shoumengUrlKeys []string = []string{"orderId", "uid", "amount",
	"coOrderId", "success"}

//手盟
type ShouMeng struct {
	MD5
}

func (this *ShouMeng) Init(params ...interface{}) (err error) {
	if err = this.MD5.Init(params...); err != nil {
		return
	}

	this.urlParamCheckKeys = &shoumengUrlKeys
	this.requireChannelUserId = true
	this.exChangeRatio = 100

	this.channelParamKeys["_gameId"] = ""
	this.channelParamKeys["_gameKey"] = ""
	this.channelParamKeys["_payKey"] = "SHOUMENG_SERCETKEY"

	this.channelTradeData = nil
	this.channelRetData = nil

	this.signHandleMethod = ""
	return
}

func (this *ShouMeng) ParseInputParam(params ...interface{}) (err error) {
	if err = this.MD5.ParseInputParam(); err != nil {
		return
	}

	orderId_key := "coOrderId"
	channelUserId_key := "uid"
	channelOrderId_key := "orderId"
	amount_key := "amount"
	discount_key := ""
	return parseTradeData_urlParam(&this.MD5.Base,
		orderId_key, channelUserId_key, channelOrderId_key,
		amount_key, discount_key)
}

func (this *ShouMeng) CheckSign(params ...interface{}) (err error) {
	format := "orderId=%s&uid=%s&amount=%s&coOrderId=%s&success=%s&secret=%s"
	this.signContent = fmt.Sprintf(format,
		this.channelOrderId,
		this.channelUserId,
		this.urlParams.Get("amount"),
		this.orderId,
		this.urlParams.Get("success"),
		this.channelParams["_payKey"])
	this.inputSign = this.urlParams.Get("sign")
	return this.MD5.CheckSign()
}

func (this *ShouMeng) CheckChannelRet(params ...interface{}) (err error) {
	tradeState := this.urlParams.Get("success") == "0"
	tradeFailDesc := `urlParam(success)!="0"`
	return this.MD5.CheckChannelRet(tradeState, tradeFailDesc)
}

func (this *ShouMeng) GetResult(params ...interface{}) (ret string) {
	format := `%s`
	succMsg := "SUCCESS"
	failMsg := "FAILURE"
	return this.MD5.GetResult(format, succMsg, failMsg)
}
