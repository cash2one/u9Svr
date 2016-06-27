package channelPayNotify

import (
	"fmt"
)

var cygameUrlKeys []string = []string{"orderid", "username", "gameid", "paytype",
	"amount", "paytime", "attach", "sign"}

type CYGame struct {
	MD5
}

func (this *CYGame) Init(params ...interface{}) (err error) {
	if err = this.MD5.Init(params...); err != nil {
		return
	}

	this.urlParamCheckKeys = &cygameUrlKeys
	this.requireChannelUserId = true
	this.exChangeRatio = 100

	this.channelParamKeys["_gameId"] = ""
	this.channelParamKeys["_gameKey"] = ""
	this.channelParamKeys["_payKey"] = "MG_APPKEY"

	this.channelTradeData = nil
	this.channelRetData = nil

	this.signHandleMethod = ""
	return
}

func (this *CYGame) ParseInputParam(params ...interface{}) (err error) {
	if err = this.MD5.ParseInputParam(); err != nil {
		return
	}

	orderId_key := "attach"
	channelUserId_key := "username"
	channelOrderId_key := "orderid"
	amount_key := "amount"
	discount_key := ""
	return parseTradeData_urlParam(&this.MD5.Base,
		orderId_key, channelUserId_key, channelOrderId_key,
		amount_key, discount_key)
}

func (this *CYGame) CheckSign(params ...interface{}) (err error) {
	format := "orderid=%s&username=%s&gameid=%s&roleid=%s&serverid=%s&paytype=%s&amount=%s&paytime=%s&attach=%s&appkey=%s"
	this.signContent = fmt.Sprintf(format,
		this.channelOrderId,
		this.urlParams.Get("username"),
		this.urlParams.Get("gameid"),
		this.urlParams.Get("roleid"),
		this.urlParams.Get("serverid"),
		this.urlParams.Get("paytype"),
		this.urlParams.Get("amount"),
		this.urlParams.Get("paytime"),
		this.urlParams.Get("attach"),
		this.channelParams["_payKey"])
	this.inputSign = this.urlParams.Get("sign")
	return this.MD5.CheckSign()
}

func (this *CYGame) GetResult(params ...interface{}) (ret string) {
	switch this.lastError {
	case err_noerror:
		ret = "success"
	case err_checkSign:
		ret = "errorSign"
	default:
		ret = "error"
	}
	this.MD5.GetResult()
	return

}
