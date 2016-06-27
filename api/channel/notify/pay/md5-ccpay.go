package channelPayNotify

import (
	"fmt"
)

var ccpayUrlKeys []string = []string{"transactionNo", "partnerTransactionNo", "statusCode",
	"productId", "orderPrice", "packageId", "sign"}

type CCPay struct {
	MD5
}

func (this *CCPay) Init(params ...interface{}) (err error) {
	if err = this.MD5.Init(params...); err != nil {
		return
	}

	this.urlParamCheckKeys = &ccpayUrlKeys
	this.requireChannelUserId = false
	this.exChangeRatio = 100

	this.channelParamKeys["_gameId"] = ""
	this.channelParamKeys["_gameKey"] = ""
	this.channelParamKeys["_payKey"] = "app_secret"

	this.channelTradeData = nil
	this.channelRetData = nil

	this.signHandleMethod = ""
	return
}

func (this *CCPay) ParseInputParam(params ...interface{}) (err error) {
	if err = this.MD5.ParseInputParam(); err != nil {
		return
	}

	orderId_key := "partnerTransactionNo"
	channelUserId_key := ""
	channelOrderId_key := "transactionNo"
	amount_key := "orderPrice"
	discount_key := ""
	return parseTradeData_urlParam(&this.MD5.Base,
		orderId_key, channelUserId_key, channelOrderId_key,
		amount_key, discount_key)
}

func (this *CCPay) CheckSign(params ...interface{}) (err error) {
	format := "orderPrice=%s&packageId=%s&partnerTransactionNo=%s&productId=%s&statusCode=%s&transactionNo=%s&%s"
	this.signContent = fmt.Sprintf(format,
		this.urlParams.Get("orderPrice"),
		this.urlParams.Get("packageId"),
		this.urlParams.Get("partnerTransactionNo"),
		this.urlParams.Get("productId"),
		this.urlParams.Get("statusCode"),
		this.urlParams.Get("transactionNo"),
		this.channelParams["_payKey"])
	this.inputSign = this.urlParams.Get("sign")
	return this.MD5.CheckSign()
}

func (this *CCPay) CheckChannelRet(params ...interface{}) (err error) {
	tradeState := this.urlParams.Get("statusCode") == "0000"
	tradeFailDesc := `urlParam(statusCode)!="0000"`
	return this.MD5.CheckChannelRet(tradeState, tradeFailDesc)
}

func (this *CCPay) GetResult(params ...interface{}) (ret string) {
	format := `%s`
	succMsg := "success"
	failMsg := "fail"
	return this.MD5.GetResult(format, succMsg, failMsg)
}
