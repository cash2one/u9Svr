package channelPayNotify

import (
	"encoding/json"
	"fmt"
	"u9/tool"
)

var xiaomiUrlKeys []string = []string{"appId", "cpOrderId", "uid", "orderId",
	"orderStatus", "payFee", "productCode", "productName", "productCount",
	"payTime", "signature"}

type Xiaomi struct {
	Base
}

func (this *Xiaomi) Init(params ...interface{}) (err error) {
	if err = this.Base.Init(params...); err != nil {
		return
	}

	this.urlParamCheckKeys = &xiaomiUrlKeys
	this.requireChannelUserId = true
	this.exChangeRatio = 1

	this.channelParamKeys["_gameId"] = "XIAOMI_APPID"
	this.channelParamKeys["_gameKey"] = ""
	this.channelParamKeys["_payKey"] = "XIAOMI_SECRETKEY"

	this.channelTradeData = nil
	this.channelRetData = nil

	return
}

func (this *Xiaomi) ParseInputParam(params ...interface{}) (err error) {
	if err = this.Base.ParseInputParam(); err != nil {
		return
	}

	orderId_key := "cpOrderId"
	channelUserId_key := "uid"
	channelOrderId_key := "orderId"
	amount_key := "payFee"
	discount_key := ""
	return parseTradeData_urlParam(&this.Base,
		orderId_key, channelUserId_key, channelOrderId_key,
		amount_key, discount_key)
}

func (this *Xiaomi) CheckSign(params ...interface{}) (err error) {
	cpUserInfo := this.urlParams.Get("cpUserInfo")
	format := `appId=%s&cpOrderId=%s&`
	if cpUserInfo != "" {
		format = format + "cpUserInfo=" + cpUserInfo + "&"
	}

	orderConsumeType := this.urlParams.Get("orderConsumeType")
	if orderConsumeType != "" {
		format = format + "orderConsumeType=" + orderConsumeType + "&"
	}

	partnerGiftConsume := this.urlParams.Get("partnerGiftConsume")
	format = format + `orderId=%s&orderStatus=%s&`
	if partnerGiftConsume != "" {
		format = format + "partnerGiftConsume=" + partnerGiftConsume + "&"
	}
	format = format + `payFee=%s&payTime=%s&productCode=%s&productCount=%s&productName=%s&uid=%s`

	content := fmt.Sprintf(format,
		this.urlParams.Get("appId"),
		this.orderId,
		this.channelOrderId,
		this.urlParams.Get("orderStatus"),
		this.urlParams.Get("payFee"),
		this.urlParams.Get("payTime"),
		this.urlParams.Get("productCode"),
		this.urlParams.Get("productCount"),
		this.urlParams.Get("productName"),
		this.channelUserId)

	payKey := this.channelParams["_payKey"]
	encryptData := string(tool.HmacSHA1Encrypt(content, payKey))
	sign := fmt.Sprintf("%x", encryptData)
	inputSign := this.urlParams.Get("signature")

	signMethod := "HmacSHA1Encrypt"
	format = "content:%s, inputSign:%s, sign:%s"
	signMsg := fmt.Sprintf(format, content, inputSign, sign)
	signState := sign == inputSign
	return this.Base.CheckSign(signState, signMethod, signMsg)
}

func (this *Xiaomi) CheckChannelRet(params ...interface{}) (err error) {
	tradeState := this.urlParams.Get("appId") == this.channelParams["_gameId"]
	tradeFailDesc := `urlParam(appId)!=channelParam(_gameId)`
	if !tradeState {
		return this.Base.CheckChannelRet(tradeState, tradeFailDesc)
	}

	tradeState = this.urlParams.Get("orderStatus") == "TRADE_SUCCESS"
	tradeFailDesc = `urlParam(orderStatus)!="TRADE_SUCCESS"`
	return this.Base.CheckChannelRet(tradeState, tradeFailDesc)
}

func (this *Xiaomi) GetResult(params ...interface{}) (ret string) {
	type XiaomiRet struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errMsg"`
	}

	xiaomiRet := new(XiaomiRet)

	switch this.lastError {
	case err_noerror:
		xiaomiRet.ErrCode = 200
		xiaomiRet.ErrMsg = "success"
	case err_orderIsNotExist:
		xiaomiRet.ErrCode = 1506
		xiaomiRet.ErrMsg = "cpOrderId error"
	case err_initChannelGameId:
		fallthrough
	case err_checkChannelRet:
		xiaomiRet.ErrCode = 1515
		xiaomiRet.ErrMsg = "appId error"
	case err_channelUserIsNotExist:
		xiaomiRet.ErrCode = 1516
		xiaomiRet.ErrMsg = "uid error"
	case err_checkSign:
		xiaomiRet.ErrCode = 1525
		xiaomiRet.ErrMsg = "signature error"
	case err_payAmountError: //订单信息不一致 error 3515
		xiaomiRet.ErrCode = 3515
		xiaomiRet.ErrMsg = "payAmount error"
	default:
		xiaomiRet.ErrCode = this.lastError
		xiaomiRet.ErrMsg = "other error"
	}
	data, _ := json.Marshal(xiaomiRet)
	ret = string(data)

	this.Base.GetResult()
	return
}
