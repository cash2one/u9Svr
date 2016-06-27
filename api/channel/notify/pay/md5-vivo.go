package channelPayNotify

import (
	"fmt"
	"strconv"
	"u9/tool"
)

type Vivo struct {
	MD5
}

func (this *Vivo) Init(params ...interface{}) (err error) {
	if err = this.MD5.Init(params...); err != nil {
		return
	}

	this.urlParamCheckKeys = &emptyUrlKeys
	this.requireChannelUserId = true
	this.exChangeRatio = 1

	this.channelParamKeys["_gameId"] = ""
	this.channelParamKeys["_gameKey"] = ""
	this.channelParamKeys["_payKey"] = "VIVO_CP_KEY"

	this.channelTradeData = nil
	this.channelRetData = nil

	this.signHandleMethod = ""
	return
}

func (this *Vivo) ParseInputParam(params ...interface{}) (err error) {
	if err = this.MD5.ParseInputParam(); err != nil {
		return
	}

	orderId_key := "cpOrderNumber"
	channelUserId_key := "uid"
	channelOrderId_key := "orderNumber"
	amount_key := "orderAmount"
	discount_key := ""
	return parseTradeData_urlParam(&this.Base,
		orderId_key, channelUserId_key, channelOrderId_key,
		amount_key, discount_key)
}

func (this *Vivo) CheckSign(params ...interface{}) (err error) {
	format := `appId=%s&cpId=%s&cpOrderNumber=%s&extInfo=%s&orderAmount=%s&orderNumber=%s&` +
		`payTime=%s&respCode=%s&respMsg=%s&tradeStatus=%s&tradeType=%s&uid=%s&%s`
	payKey := this.channelParams["_payKey"]
	this.signContent = fmt.Sprintf(format,
		this.urlParams.Get("appId"),
		this.urlParams.Get("cpId"),
		this.urlParams.Get("cpOrderNumber"),
		this.urlParams.Get("extInfo"),
		this.urlParams.Get("orderAmount"),
		this.urlParams.Get("orderNumber"),
		this.urlParams.Get("payTime"),
		this.urlParams.Get("respCode"),
		this.urlParams.Get("respMsg"),
		this.urlParams.Get("tradeStatus"),
		this.urlParams.Get("tradeType"),
		this.urlParams.Get("uid"),
		tool.Md5([]byte(payKey)))
	this.inputSign = this.urlParams.Get("signature")
	return this.MD5.CheckSign()
}

func (this *Vivo) GetResult(params ...interface{}) (ret string) {
	if this.lastError != err_noerror {
		//this.ctx.Abort(403, msg)
		msg := "lastError:" + strconv.Itoa(this.lastError)
		this.ctx.ResponseWriter.WriteHeader(403)
		this.ctx.ResponseWriter.Write([]byte(msg))
		//panic(msg)
	}

	this.MD5.GetResult()
	return
}
