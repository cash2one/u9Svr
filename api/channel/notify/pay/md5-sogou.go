package channelPayNotify

import (
	"fmt"
)

var sogouUrlKeys []string = []string{"gid", "sid", "uid", "oid", "date", "amount1",
	"realAmount", "auth"}

//当乐
type Sogou struct {
	MD5
}

func (this *Sogou) Init(params ...interface{}) (err error) {
	if err = this.MD5.Init(params...); err != nil {
		return
	}

	this.urlParamCheckKeys = &sogouUrlKeys
	this.requireChannelUserId = true
	this.exChangeRatio = 100

	this.channelParamKeys["_gameId"] = "SOGOU_GAMEID"
	this.channelParamKeys["_gameKey"] = ""
	this.channelParamKeys["_payKey"] = "SOGOU_PAYKEY"

	this.channelTradeData = nil
	this.channelRetData = nil

	this.signHandleMethod = ""
	return
}

func (this *Sogou) ParseInputParam(params ...interface{}) (err error) {
	if err = this.MD5.ParseInputParam(); err != nil {
		return
	}

	orderId_key := "appdata"
	channelUserId_key := "uid"
	channelOrderId_key := "oid"
	amount_key := "realAmount"
	discount_key := ""
	return parseTradeData_urlParam(&this.Base,
		orderId_key, channelUserId_key, channelOrderId_key,
		amount_key, discount_key)
}

func (this *Sogou) CheckSign(params ...interface{}) (err error) {
	format := "amount1=%s&amount2=%s&appdata=%s&date=%s&gid=%s&oid=%s&realAmount=%s&role=%s&sid=%s&time=%s&uid=%s&%s"
	this.signContent = fmt.Sprintf(format,
		this.urlParams.Get("amount1"),
		this.urlParams.Get("amount2"),
		this.orderId,
		this.urlParams.Get("date"),
		this.channelParams["_gameId"],
		this.channelOrderId,
		this.urlParams.Get("realAmount"),
		this.urlParams.Get("role"),
		this.urlParams.Get("sid"),
		this.urlParams.Get("time"),
		this.channelUserId,
		this.channelParams["_payKey"])
	this.inputSign = this.urlParams.Get("auth")
	return this.MD5.CheckSign()
}

func (this *Sogou) GetResult(params ...interface{}) (ret string) {
	switch this.lastError {
	case err_noerror:
		ret = "OK"
	case err_checkSign:
		ret = "ERR_200"
	default:
		ret = "ERR_100"
	}
	this.MD5.GetResult()
	return
}
