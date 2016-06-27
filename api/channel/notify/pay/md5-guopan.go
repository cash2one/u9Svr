package channelPayNotify

import ()

var guopanUrlKeys []string = []string{"trade_no", "serialNumber",
	"money", "status", "t", "sign", "reserved"}

type Guopan struct {
	MD5
}

func (this *Guopan) Init(params ...interface{}) (err error) {
	if err = this.Base.Init(params...); err != nil {
		return
	}

	this.urlParamCheckKeys = &guopanUrlKeys
	this.requireChannelUserId = false
	this.exChangeRatio = 100

	this.channelParamKeys["_gameId"] = ""
	this.channelParamKeys["_gameKey"] = ""
	this.channelParamKeys["_payKey"] = "GUOPAN_SERVER_SECRETKEY"

	this.channelTradeData = nil
	this.channelRetData = nil

	this.signHandleMethod = ""
	return
}

func (this *Guopan) ParseInputParam(params ...interface{}) (err error) {
	if err = this.MD5.ParseInputParam(); err != nil {
		return
	}

	orderId_key := "serialNumber"
	channelUserId_key := ""
	channelOrderId_key := "trade_no"
	amount_key := "money"
	discount_key := ""
	return parseTradeData_urlParam(&this.MD5.Base,
		orderId_key, channelUserId_key, channelOrderId_key,
		amount_key, discount_key)
}

func (this *Guopan) CheckSign(params ...interface{}) (err error) {
	this.signContent = this.urlParams.Get("serialNumber") +
		this.urlParams.Get("money") +
		this.urlParams.Get("status") +
		this.urlParams.Get("t") +
		this.channelParams["_payKey"]
	this.inputSign = this.urlParams.Get("Sign")
	return this.MD5.CheckSign()
}

func (this *Guopan) CheckChannelRet(params ...interface{}) (err error) {
	tradeState := this.urlParams.Get("status") == "1"
	tradeFailDesc := `urlParam(status)!="1"`
	return this.MD5.CheckChannelRet(tradeState, tradeFailDesc)
}

func (this *Guopan) GetResult(params ...interface{}) (ret string) {
	format := `%s`
	succMsg := "success"
	failMsg := "fail"
	return this.MD5.GetResult(format, succMsg, failMsg)
}
