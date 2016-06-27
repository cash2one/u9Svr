package channelPayNotify

import (
	"fmt"
)

var amigoUrlKeys []string = []string{"api_key", "close_time", "create_time", "deal_price", "out_order_no",
	"pay_channel", "submit_time", "user_id", "sign"}

type Amigo struct {
	Rsa
}

func (this *Amigo) Init(params ...interface{}) (err error) {
	if err = this.Rsa.Init(params...); err != nil {
		return
	}

	this.urlParamCheckKeys = &amigoUrlKeys
	this.requireChannelUserId = false
	this.exChangeRatio = 100

	this.channelParamKeys["_gameId"] = ""
	this.channelParamKeys["_gameKey"] = "AMIGO_APIKEY"
	this.channelParamKeys["_payKey"] = "AMIGO_PUBLICKEY"

	this.channelTradeData = nil
	this.channelRetData = nil

	this.signMode = 0
	return
}

func (this *Amigo) ParseInputParam(params ...interface{}) (err error) {
	if err = this.Rsa.ParseInputParam(); err != nil {
		return
	}

	orderId_key := "out_order_no"
	channelUserId_key := ""
	channelOrderId_key := ""
	amount_key := "deal_price"
	discount_key := ""
	return parseTradeData_urlParam(&this.Rsa.Base,
		orderId_key, channelUserId_key, channelOrderId_key,
		amount_key, discount_key)
}

func (this *Amigo) CheckSign(params ...interface{}) (err error) {
	format := "api_key=%s&close_time=%s&create_time=%s&deal_price=%s&out_order_no=%s&pay_channel=%s&submit_time=%s&user_id=%s"
	this.signContent = fmt.Sprintf(format,
		this.channelParams["_gameKey"],
		this.urlParams.Get("close_time"),
		this.urlParams.Get("create_time"),
		this.urlParams.Get("deal_price"),
		this.orderId,
		this.urlParams.Get("pay_channel"),
		this.urlParams.Get("submit_time"),
		"null")
	this.inputSign = this.urlParams.Get("sign")
	return this.Rsa.CheckSign()
}

func (this *Amigo) GetResult(params ...interface{}) (ret string) {
	format := `%s`
	succMsg := "success"
	failMsg := "failure"
	return this.Rsa.GetResult(format, succMsg, failMsg)
}
