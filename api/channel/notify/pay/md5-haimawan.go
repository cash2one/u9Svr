package channelPayNotify

import (
	"fmt"
	"net/url"
)

var haimawanUrlKeys []string = []string{"notify_time", "appid", "out_trade_no", "total_fee",
	"subject", "body", "trade_status", "sign"}

//海马玩
type Haimawan struct {
	MD5
}

func (this *Haimawan) Init(params ...interface{}) (err error) {
	if err = this.MD5.Init(params...); err != nil {
		return
	}

	this.urlParamCheckKeys = &haimawanUrlKeys
	this.requireChannelUserId = false
	this.exChangeRatio = 100

	this.channelParamKeys["_gameId"] = ""
	this.channelParamKeys["_gameKey"] = ""
	this.channelParamKeys["_payKey"] = "HM_SERVERKEY"

	this.channelTradeData = nil
	this.channelRetData = nil

	this.signHandleMethod = ""
	return
}

func (this *Haimawan) ParseInputParam(params ...interface{}) (err error) {
	if err = this.MD5.ParseInputParam(); err != nil {
		return
	}

	orderId_key := "out_trade_no"
	channelUserId_key := ""
	channelOrderId_key := ""
	amount_key := "total_fee"
	discount_key := ""
	return parseTradeData_urlParam(&this.MD5.Base,
		orderId_key, channelUserId_key, channelOrderId_key,
		amount_key, discount_key)
}

func (this *Haimawan) CheckSign(params ...interface{}) (err error) {
	format := "notify_time=%s&appid=%s&out_trade_no=%s&total_fee=%s&subject=%s&body=%s&trade_status=%s%s"
	this.signContent = fmt.Sprintf(format,
		url.QueryEscape(this.urlParams.Get("notify_time")),
		url.QueryEscape(this.urlParams.Get("appid")),
		url.QueryEscape(this.urlParams.Get("out_trade_no")),
		url.QueryEscape(this.urlParams.Get("total_fee")),
		url.QueryEscape(this.urlParams.Get("subject")),
		url.QueryEscape(this.urlParams.Get("body")),
		url.QueryEscape(this.urlParams.Get("trade_status")),
		this.channelParams["_payKey"])

	this.inputSign = this.urlParams.Get("sign")
	return this.MD5.CheckSign()
}

func (this *Haimawan) CheckChannelRet(params ...interface{}) (err error) {
	tradeState := this.urlParams.Get("trade_status") == "1"
	tradeFailDesc := `urlParam(trade_status)!="1"`
	return this.MD5.CheckChannelRet(tradeState, tradeFailDesc)
}

func (this *Haimawan) GetResult(params ...interface{}) (ret string) {
	format := `%s`
	succMsg := "success"
	failMsg := "failure"
	return this.MD5.GetResult(format, succMsg, failMsg)
}
