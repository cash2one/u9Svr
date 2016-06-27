package channelPayNotify

import (
	"fmt"
)

var mzwUrlKeys []string = []string{"appkey", "orderID", "money", "uid",
	"extern", "sign"}

//拇指玩
type MZW struct {
	MD5
}

func (this *MZW) Init(params ...interface{}) (err error) {
	if err = this.MD5.Init(params...); err != nil {
		return
	}

	this.urlParamCheckKeys = &mzwUrlKeys
	this.requireChannelUserId = true
	this.exChangeRatio = 100

	this.channelParamKeys["_gameId"] = ""
	this.channelParamKeys["_gameKey"] = "MZWAPPKEY"
	this.channelParamKeys["_payKey"] = "MZWPAYKEY"

	this.channelTradeData = nil
	this.channelRetData = nil

	this.signHandleMethod = ""
	return
}

func (this *MZW) ParseInputParam(params ...interface{}) (err error) {
	if err = this.MD5.ParseInputParam(); err != nil {
		return
	}

	orderId_key := "extern"
	channelUserId_key := "uid"
	channelOrderId_key := "orderID"
	amount_key := "money"
	discount_key := ""
	return parseTradeData_urlParam(&this.Base,
		orderId_key, channelUserId_key, channelOrderId_key,
		amount_key, discount_key)
}

func (this *MZW) CheckSign(params ...interface{}) (err error) {
	format := "%s%s%s%s%s%s%s%s%s"
	this.signContent = fmt.Sprintf(format,
		this.channelParams["_gameKey"],
		this.channelOrderId,
		this.urlParams.Get("productName"),
		this.urlParams.Get("productDesc"),
		this.urlParams.Get("productID"),
		this.urlParams.Get("money"),
		this.channelUserId,
		this.urlParams.Get("extern"),
		this.channelParams["_payKey"])
	this.inputSign = this.urlParams.Get("sign")
	return this.MD5.CheckSign()
}

func (this *MZW) GetResult(params ...interface{}) (ret string) {
	format := `%s`
	succMsg := "SUCCESS"
	failMsg := "FAILURE"
	return this.MD5.GetResult(format, succMsg, failMsg)
}
