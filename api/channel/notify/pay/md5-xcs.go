package channelPayNotify

import (
	"fmt"
)

var caishenUrlKeys []string = []string{"product_id", "order_id", "price", "game_uid", "u_id", "xcs_order", "game_id", "sign"}

//小财神
type Xcs struct {
	MD5
}

func (this *Xcs) Init(params ...interface{}) (err error) {
	if err = this.MD5.Init(params...); err != nil {
		return
	}

	this.urlParamCheckKeys = &caishenUrlKeys
	this.requireChannelUserId = true
	this.exChangeRatio = 100

	this.channelParamKeys["_gameId"] = ""
	this.channelParamKeys["_gameKey"] = ""
	this.channelParamKeys["_payKey"] = "PAYKEY"

	this.channelTradeData = nil
	this.channelRetData = nil

	this.signHandleMethod = ""
	return
}

func (this *Xcs) ParseInputParam(params ...interface{}) (err error) {
	if err = this.MD5.ParseInputParam(); err != nil {
		return
	}

	orderId_key := "order_id"
	channelUserId_key := "u_id"
	channelOrderId_key := "xcs_order"
	amount_key := "price"
	discount_key := ""
	return parseTradeData_urlParam(&this.MD5.Base,
		orderId_key, channelUserId_key, channelOrderId_key,
		amount_key, discount_key)
}

func (this *Xcs) CheckSign(params ...interface{}) (err error) {
	this.signContent = fmt.Sprintf("%s_%s_%s_%s_%s_%s_%s_%s",
		this.orderId,
		this.urlParams.Get("product_id"),
		this.urlParams.Get("price"),
		this.urlParams.Get("game_uid"),
		this.channelUserId,
		this.channelOrderId,
		this.urlParams.Get("game_id"),
		this.channelParams["_payKey"])
	this.inputSign = this.urlParams.Get("sign")
	return this.MD5.CheckSign()
}

func (this *Xcs) GetResult(params ...interface{}) (ret string) {
	switch this.lastError {
	case err_noerror:
		ret = "0000"
	case err_checkSign:
		ret = "0002"
	default:
		ret = "unknow"
	}
	this.MD5.GetResult()
	return
}
