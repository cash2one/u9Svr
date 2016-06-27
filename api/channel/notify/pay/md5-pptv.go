package channelPayNotify

import (
	"fmt"
)

var pptvUrlKeys []string = []string{"username", "oid", "amount", "extra", "time", "sign"}

//PPTV
type PPTV struct {
	MD5
}

func (this *PPTV) Init(params ...interface{}) (err error) {
	if err = this.MD5.Init(params...); err != nil {
		return
	}

	this.urlParamCheckKeys = &pptvUrlKeys
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

func (this *PPTV) ParseInputParam(params ...interface{}) (err error) {
	if err = this.MD5.ParseInputParam(); err != nil {
		return
	}

	orderId_key := "extra"
	channelUserId_key := "username"
	channelOrderId_key := "oid"
	amount_key := "amount"
	discount_key := ""
	return parseTradeData_urlParam(&this.MD5.Base,
		orderId_key, channelUserId_key, channelOrderId_key,
		amount_key, discount_key)
}

func (this *PPTV) CheckSign(params ...interface{}) (err error) {
	format := "%s%s%s%s%s%s%s"
	this.signContent = fmt.Sprintf(format,
		this.urlParams.Get("sid"),
		this.channelUserId,
		this.urlParams.Get("roid"),
		this.channelOrderId,
		this.urlParams.Get("amount"),
		this.urlParams.Get("time"),
		this.channelParams["_payKey"])
	this.inputSign = this.urlParams.Get("sign")
	return this.MD5.CheckSign()
}

func (this *PPTV) GetResult(params ...interface{}) (ret string) {

	format := `{"code":"%d","message":"%s"}`
	switch this.lastError {
	case err_noerror:
		ret = fmt.Sprintf(format, 1, "success")
	case err_orderIsNotExist:
		fallthrough
	case err_prepareOrderRequest:
		ret = fmt.Sprintf(format, 3, "order_fail")
	case err_checkSign:
		ret = fmt.Sprintf(format, 4, "sign_fail")
	case err_channelUserIsNotExist:
		ret = fmt.Sprintf(format, 5, "user_fail")
	case err_payAmountError:
		ret = fmt.Sprintf(format, 6, "money_fail")
	case err_tradeFail:
		fallthrough
	default:
		ret = fmt.Sprintf(format, 0, "pay_fail")
	}

	this.MD5.GetResult()
	return
}
