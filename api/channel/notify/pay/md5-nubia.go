package channelPayNotify

import (
	"fmt"
)

var nubiaUrlKeys []string = []string{"order_no", "data_timestamp", "pay_success", "order_sign","uid","amount"}

//努比亚
type Nubia struct {
	MD5
}

func (this *Nubia) Init(params ...interface{}) (err error) {
	if err = this.MD5.Init(params...); err != nil {
		return
	}

	this.urlParamCheckKeys = &nubiaUrlKeys
	this.requireChannelUserId = true
	this.exChangeRatio = 100

	this.channelParamKeys["_gameId"] = "NUBIA_APPID"
	this.channelParamKeys["_gameKey"] = ""
	this.channelParamKeys["_payKey"] = "NUBIA_APPSECRET_KEY"

	this.channelTradeData = nil
	this.channelRetData = nil

	this.signHandleMethod = ""
	return
}

func (this *Nubia) ParseInputParam(params ...interface{}) (err error) {
	if err = this.MD5.ParseInputParam(); err != nil {
		return
	}

	orderId_key := "order_no"
	channelUserId_key := "uid"
	channelOrderId_key := ""
	amount_key := "amount"
	discount_key := ""
	return parseTradeData_urlParam(&this.Base,
		orderId_key, channelUserId_key, channelOrderId_key,
		amount_key, discount_key)
}

func (this *Nubia) CheckSign(params ...interface{}) (err error) {
	this.inputSign = this.urlParams.Get("order_sign")

		format := "amount=%s&app_id=%s&data_timestamp=%s&number=1&order_no=%s&pay_success=%s&product_des=%s&product_name=%s&uid=%s:%s:%s"
		this.signContent = fmt.Sprintf(format,
			this.urlParams.Get("amount"),
			this.channelParams["_gameId"],
			this.urlParams.Get("data_timestamp"),
			this.orderId,
			this.urlParams.Get("pay_success"),
			this.urlParams.Get("product_des"),
			this.urlParams.Get("product_name"),
			this.urlParams.Get("uid"),
			this.channelParams["_gameId"],
			this.channelParams["_payKey"])
	return this.MD5.CheckSign()
}

func (this *Nubia) CheckChannelRet(params ...interface{}) (err error) {
	tradeState := this.urlParams.Get("pay_success") == "1"
	tradeFailDesc := `urlParam(pay_success)!=1`
	return this.MD5.CheckChannelRet(tradeState, tradeFailDesc)
}

func (this *Nubia) GetResult(params ...interface{}) (ret string) {
	format := `{"code":%s,"data":{},"message":"%s"}`
	switch this.lastError {
	case err_noerror:
		ret = fmt.Sprintf(format, "0", "成功")
	case err_checkSign:
		ret = fmt.Sprintf(format, "90000", "Sign无效")
	case err_initChannelGameId:
		fallthrough
	case err_initChannelPayKey:
		ret = fmt.Sprintf(format, "90000", "参数无效")
	default:
		ret = fmt.Sprintf(format, "10000", "接受失败")
	}
		this.MD5.GetResult()
	return 
}
