package channelPayNotify

import (
	"fmt"
)

var paojiaoUrlKeys []string = []string{"uid", "orderNo", "price", "status", "gameId", "payTime", "ext","sign"}

//泡椒
type PaoJiao struct {
	MD5
}

func (this *PaoJiao) Init(params ...interface{}) (err error) {
	if err = this.MD5.Init(params...); err != nil {
		return
	}

	this.urlParamCheckKeys = &paojiaoUrlKeys
	this.requireChannelUserId = true
	this.exChangeRatio = 100

	this.channelParamKeys["_gameId"] = ""
	this.channelParamKeys["_gameKey"] = ""
	this.channelParamKeys["_payKey"] = "PAOJIAO_SERVERSECRET"

	this.channelTradeData = nil
	this.channelRetData = nil

	this.signHandleMethod = ""
	return
}

func (this *PaoJiao) ParseInputParam(params ...interface{}) (err error) {
	if err = this.MD5.ParseInputParam(); err != nil {
		return
	}

	orderId_key := "ext"
	channelUserId_key := "uid"
	channelOrderId_key := "orderNo"
	amount_key := "price"
	discount_key := ""
	return parseTradeData_urlParam(&this.MD5.Base,
		orderId_key, channelUserId_key, channelOrderId_key,
		amount_key, discount_key)
}

func (this *PaoJiao) CheckSign(params ...interface{}) (err error) {
	format := "uid=%sprice=%sorderNo=%sremark=%sstatus=%ssubject=%sgameId=%spayTime=%sext=%s%s"
	this.signContent = fmt.Sprintf(format,
		this.channelUserId,
		this.urlParams.Get("price"),
		this.channelOrderId,
		this.urlParams.Get("remark"),
		this.urlParams.Get("status"),
		this.urlParams.Get("subject"),
		this.urlParams.Get("gameId"),
		this.urlParams.Get("payTime"),
		this.urlParams.Get("ext"),
		this.channelParams["_payKey"])
	this.inputSign = this.urlParams.Get("sign")
	return this.MD5.CheckSign()
}

func (this *PaoJiao) CheckChannelRet(params ...interface{}) (err error) {
	tradeState := this.urlParams.Get("status") == "5"
	tradeFailDesc := `urlParam(status)!="5"`
	return this.MD5.CheckChannelRet(tradeState, tradeFailDesc)
}

func (this *PaoJiao) GetResult(params ...interface{}) (ret string) {
	format := `%s`
	succMsg := "success"
	failMsg := "failure"
	return this.MD5.GetResult(format, succMsg, failMsg)
}
