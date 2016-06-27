package channelPayNotify

import (
	"fmt"
	"net/url"
)

var mangoUrlKeys []string = []string{"ret_code", "error_msg", "aid", "order_no"}

//芒果玩
type Mango struct {
	MD5
}

func (this *Mango) Init(params ...interface{}) (err error) {
	if err = this.MD5.Init(params...); err != nil {
		return
	}

	this.urlParamCheckKeys = &mangoUrlKeys
	this.requireChannelUserId = true
	this.exChangeRatio = 100

	this.channelParamKeys["_gameId"] = ""
	this.channelParamKeys["_gameKey"] = ""
	this.channelParamKeys["_payKey"] = "SKEY"

	this.channelTradeData = nil
	this.channelRetData = nil

	this.signHandleMethod = ""
	return
}

func (this *Mango) ParseInputParam(params ...interface{}) (err error) {
	if err = this.MD5.ParseInputParam(); err != nil {
		return
	}

	orderId_key := "extension_field"
	channelUserId_key := "aid"
	channelOrderId_key := "order_no"
	amount_key := "money"
	discount_key := ""
	return parseTradeData_urlParam(&this.Base,
		orderId_key, channelUserId_key, channelOrderId_key,
		amount_key, discount_key)
}

func (this *Mango) CheckSign(params ...interface{}) (err error) {
	format := "%s%s%s%s%s%s%s%s%s%s"
	this.signContent = fmt.Sprintf(format,
		this.urlParams.Get("ret_code"),
		this.urlParams.Get("aid"),
		this.urlParams.Get("gid"),
		this.urlParams.Get("cid"),
		this.channelParams["_payKey"],
		this.urlParams.Get("ts"),
		url.QueryEscape(this.urlParams.Get("order_no")),
		url.QueryEscape(this.urlParams.Get("pay_order")),
		this.urlParams.Get("money"),
		this.urlParams.Get("pay_type"))
	this.inputSign = this.urlParams.Get("enc")
	return this.MD5.CheckSign()
}

func (this *Mango) CheckChannelRet(params ...interface{}) (err error) {
	tradeState := this.urlParams.Get("ret_code") == "0"
	tradeFailDesc := `urlParam(ret_code)!="0"`
	return this.MD5.CheckChannelRet(tradeState, tradeFailDesc)
}

func (this *Mango) GetResult(params ...interface{}) (ret string) {
	format := `%s`
	succMsg := "success"
	failMsg := "failure"
	return this.MD5.GetResult(format, succMsg, failMsg)
}
