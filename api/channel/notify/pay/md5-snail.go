package channelPayNotify

import (
	"fmt"
)

var snailUrlKeys []string = []string{"AppId", "ConsumeStreamId", "CooOrderSerial",
	"Uin", "OriginalMoney", "OrderMoney", "PayStatus"}

//免商店
type Snail struct {
	MD5
}

func (this *Snail) Init(params ...interface{}) (err error) {
	if err = this.MD5.Init(params...); err != nil {
		return
	}

	this.urlParamCheckKeys = &snailUrlKeys
	this.requireChannelUserId = true
	this.exChangeRatio = 100

	this.channelParamKeys["_gameId"] = "SNAIL_APPID"
	this.channelParamKeys["_gameKey"] = ""
	this.channelParamKeys["_payKey"] = "SNAIL_APPKEY"

	this.channelTradeData = nil
	this.channelRetData = new(m4399ChannelRet)

	this.signHandleMethod = ""
	return
}

func (this *Snail) ParseInputParam(params ...interface{}) (err error) {
	if err = this.MD5.ParseInputParam(); err != nil {
		return
	}

	orderId_key := "CooOrderSerial"
	channelUserId_key := "Uin"
	channelOrderId_key := "ConsumeStreamId"
	amount_key := "OriginalMoney"
	discount_key := ""
	return parseTradeData_urlParam(&this.Base,
		orderId_key, channelUserId_key, channelOrderId_key,
		amount_key, discount_key)
}

func (this *Snail) CheckSign(params ...interface{}) (err error) {
	format := "%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s"
	this.signContent = fmt.Sprintf(format,
		this.channelParams["_gameId"],
		this.urlParams.Get("Act"),
		this.urlParams.Get("ProductName"),
		this.channelOrderId,
		this.orderId, this.channelUserId,
		this.urlParams.Get("GoodsId"),
		this.urlParams.Get("GoodsInfo"),
		this.urlParams.Get("GoodsCount"),
		this.urlParams.Get("OriginalMoney"),
		this.urlParams.Get("OrderMoney"),
		this.urlParams.Get("Note"),
		this.urlParams.Get("PayStatus"),
		this.urlParams.Get("CreateTime"),
		this.channelParams["_payKey"])
	this.inputSign = this.urlParams.Get("Sign")
	return this.MD5.CheckSign()
}

func (this *Snail) CheckChannelRet(params ...interface{}) (err error) {
	tradeState := this.urlParams.Get("PayStatus") == "1"
	tradeFailDesc := `urlParam(PayStatus)!="1"`
	return this.MD5.CheckChannelRet(tradeState, tradeFailDesc)
}

func (this *Snail) GetResult(params ...interface{}) (ret string) {
	format := `{"ErrorCode":"%s","ErrprDesc":"%s"}`
	switch this.lastError {
	case err_noerror:
		ret = fmt.Sprintf(format, "1", "接受成功")
	case err_checkSign:
		ret = fmt.Sprintf(format, "5", "Sign无效")
	case err_initChannelGameId:
		fallthrough
	case err_initChannelPayKey:
		ret = fmt.Sprintf(format, "2", "AppId无效")
	default:
		ret = fmt.Sprintf(format, "0", "接受失败")
	}

	this.MD5.GetResult()
	return
}
