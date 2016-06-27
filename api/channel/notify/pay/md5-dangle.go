package channelPayNotify

import (
	"fmt"
)

var dangleUrlKeys []string = []string{"result", "money", "order", "mid", "time", "ext", "signature"}

//当乐
type Dangle struct {
	MD5
}

func (this *Dangle) Init(params ...interface{}) (err error) {
	if err = this.MD5.Init(params...); err != nil {
		return
	}

	this.urlParamCheckKeys = &dangleUrlKeys
	this.requireChannelUserId = true
	this.exChangeRatio = 100

	this.channelParamKeys["_gameId"] = ""
	this.channelParamKeys["_gameKey"] = ""
	this.channelParamKeys["_payKey"] = "DANGLE_PAYMENT_KEY"

	this.channelTradeData = nil
	this.channelRetData = nil

	this.signHandleMethod = ""
	return
}

func (this *Dangle) ParseInputParam(params ...interface{}) (err error) {
	if err = this.MD5.ParseInputParam(); err != nil {
		return
	}

	orderId_key := "ext"
	channelUserId_key := "mid"
	channelOrderId_key := "order"
	amount_key := "money"
	discount_key := ""
	return parseTradeData_urlParam(&this.MD5.Base,
		orderId_key, channelUserId_key, channelOrderId_key,
		amount_key, discount_key)
}

func (this *Dangle) CheckSign(params ...interface{}) (err error) {
	format := "order=%s&money=%s&mid=%s&time=%s&result=%s&ext=%s&key=%s"
	this.signContent = fmt.Sprintf(format,
		this.channelOrderId,
		this.urlParams.Get("money"),
		this.channelUserId,
		this.urlParams.Get("time"),
		this.urlParams.Get("result"),
		this.urlParams.Get("ext"),
		this.channelParams["_payKey"])
	this.inputSign = this.urlParams.Get("signature")
	return this.MD5.CheckSign()
}

func (this *Dangle) CheckChannelRet(params ...interface{}) (err error) {
	tradeState := this.urlParams.Get("result") == "1"
	tradeFailDesc := `urlParam(result)!="1"`
	return this.MD5.CheckChannelRet(tradeState, tradeFailDesc)
}

func (this *Dangle) GetResult(params ...interface{}) (ret string) {
	format := `%s`
	succMsg := "success"
	failMsg := "failure"
	return this.MD5.GetResult(format, succMsg, failMsg)
}
