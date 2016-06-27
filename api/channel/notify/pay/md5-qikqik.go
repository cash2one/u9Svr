package channelPayNotify

import (
	"fmt"
)

var qikqikUrlKeys []string = []string{"uid", "cporder", "money", "order", "cpappid"}

//7k7k
type QikQik struct {
	MD5
}

func (this *QikQik) Init(params ...interface{}) (err error) {
	if err = this.MD5.Init(params...); err != nil {
		return
	}

	this.urlParamCheckKeys = &qikqikUrlKeys
	this.requireChannelUserId = true
	this.exChangeRatio = 100

	this.channelParamKeys["_gameId"] = "APPID"
	this.channelParamKeys["_gameKey"] = ""
	this.channelParamKeys["_payKey"] = "APP_SECRET"

	this.channelTradeData = nil
	this.channelRetData = nil

	this.signHandleMethod = ""
	return
}

func (this *QikQik) ParseInputParam(params ...interface{}) (err error) {
	if err = this.MD5.ParseInputParam(); err != nil {
		return
	}

	orderId_key := "cporder"
	channelUserId_key := "uid"
	channelOrderId_key := "order"
	amount_key := "money"
	discount_key := ""
	return parseTradeData_urlParam(&this.Base,
		orderId_key, channelUserId_key, channelOrderId_key,
		amount_key, discount_key)
}

func (this *QikQik) CheckSign(params ...interface{}) (err error) {

	format := "%s%s%s%s%s%s"
	this.signContent = fmt.Sprintf(format,
		this.channelUserId,
		this.orderId,
		this.urlParams.Get("money"),
		this.channelOrderId,
		this.channelParams["_gameId"],
		this.channelParams["_payKey"])
	this.inputSign = this.urlParams.Get("sign")
	return this.MD5.CheckSign()
}

func (this *QikQik) GetResult(params ...interface{}) (ret string) {
	format := `%s`
	succMsg := "success"
	failMsg := "fail"
	return this.MD5.GetResult(format, succMsg, failMsg)
}
