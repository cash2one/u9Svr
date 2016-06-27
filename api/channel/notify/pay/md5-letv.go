package channelPayNotify

import (
	"fmt"
)

var letvUrlKeys []string = []string{"app_id", "lepay_order_no", "letv_user_id", "out_trade_no",
	"pay_time", "price", "product_id", "sign", "sign_type", "trade_result", "version"}

//乐视
type LeTV struct {
	MD5
}

func (this *LeTV) Init(params ...interface{}) (err error) {
	if err = this.MD5.Init(params...); err != nil {
		return
	}

	this.urlParamCheckKeys = &letvUrlKeys
	this.requireChannelUserId = true
	this.exChangeRatio = 100

	this.channelParamKeys["_gameId"] = "lepay_appid"
	this.channelParamKeys["_gameKey"] = ""
	this.channelParamKeys["_payKey"] = "lepay_secretkey"

	this.channelTradeData = nil
	this.channelRetData = nil

	this.signHandleMethod = ""
	return
}

func (this *LeTV) ParseInputParam(params ...interface{}) (err error) {
	if err = this.MD5.ParseInputParam(); err != nil {
		return
	}

	orderId_key := "cooperator_order_no"
	channelUserId_key := "letv_user_id"
	channelOrderId_key := "out_trade_no"
	amount_key := "price"
	discount_key := ""
	return parseTradeData_urlParam(&this.Base,
		orderId_key, channelUserId_key, channelOrderId_key,
		amount_key, discount_key)
}

func (this *LeTV) CheckSign(params ...interface{}) (err error) {
	channelApiVersion := this.channelParams["_version"]
	this.inputSign = this.urlParams.Get("sign")
	switch channelApiVersion {
	case "1.2.2":
		format := "app_id=%s&lepay_order_no=%s&letv_user_id=%s&out_trade_no=%s&pay_time=%s&price=%s&product_id=%s&sign_type=%s&trade_result=%s&version=%s&key=%s"
		this.signContent = fmt.Sprintf(format,
			this.channelParams["_gameId"],
			this.urlParams.Get("lepay_order_no"),
			this.channelUserId,
			this.channelOrderId,
			this.urlParams.Get("pay_time"),
			this.urlParams.Get("price"),
			this.urlParams.Get("product_id"),
			"MD5",
			this.urlParams.Get("trade_result"),
			this.urlParams.Get("version"),
			this.channelParams["_payKey"])
	case "2.2.1":
		fallthrough
	default:
		format := "app_id=%s&cooperator_order_no=%s&extra_info=%s&lepay_order_no=%s&letv_user_id=%s&original_price=%s&out_trade_no=%s&pay_time=%s&price=%s&product_id=%s&sign_type=%s&trade_result=%s&version=%s&key=%s"
		this.signContent = fmt.Sprintf(format,
			this.channelParams["_gameId"],
			this.orderId,
			this.urlParams.Get("extra_info"),
			this.urlParams.Get("lepay_order_no"),
			this.channelUserId,
			this.urlParams.Get("original_price"),
			this.urlParams.Get("out_trade_no"),
			this.urlParams.Get("pay_time"),
			this.urlParams.Get("price"),
			this.urlParams.Get("product_id"),
			this.urlParams.Get("sign_type"),
			this.urlParams.Get("trade_result"),
			this.urlParams.Get("version"),
			this.channelParams["_payKey"])
	}
	return this.MD5.CheckSign()
}

func (this *LeTV) CheckChannelRet(params ...interface{}) (err error) {
	tradeState := this.urlParams.Get("trade_result") == "TRADE_SUCCESS"
	tradeFailDesc := `urlParam(trade_result)!="TRADE_SUCCESS"`
	return this.MD5.CheckChannelRet(tradeState, tradeFailDesc)
}

func (this *LeTV) GetResult(params ...interface{}) (ret string) {
	format := `%s`
	succMsg := "success"
	failMsg := "fail"
	return this.MD5.GetResult(format, succMsg, failMsg)
}
