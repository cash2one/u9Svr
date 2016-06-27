package channelPayNotify

import (
	"encoding/json"
	"fmt"
)

var meizuUrlKeys []string = []string{"notify_time", "notify_id", "order_id",
	"app_id", "uid", "partner_id", "cp_order_id", "product_id", "total_price",
	"trade_status", "create_time", "pay_time", "sign", "sign_type"}

type meizuChannelRet struct {
	Code     string `json:"code"` //200 成功发货；120013 尚未发货；120014 发货失败; 900000 未知异常
	Message  string `json:"message"`
	Value    string `json:"value"`
	Redirect string `json:"redirect"`
}

type Meizu struct {
	MD5
}

func (this *Meizu) Init(params ...interface{}) (err error) {
	if err = this.MD5.Init(params...); err != nil {
		return
	}

	this.urlParamCheckKeys = &meizuUrlKeys
	this.requireChannelUserId = true
	this.exChangeRatio = 100

	this.channelParamKeys["_gameId"] = ""
	this.channelParamKeys["_gameKey"] = ""
	this.channelParamKeys["_payKey"] = "MEIZU_APPSECRET"

	this.channelTradeData = nil
	this.channelRetData = new(meizuChannelRet)

	this.signHandleMethod = ""
	return
}

func (this *Meizu) ParseInputParam(params ...interface{}) (err error) {
	if err = this.MD5.ParseInputParam(); err != nil {
		return
	}

	orderId_key := "cp_order_id"
	channelUserId_key := "uid"
	channelOrderId_key := "order_id"
	amount_key := "total_price"
	discount_key := ""
	return parseTradeData_urlParam(&this.Base,
		orderId_key, channelUserId_key, channelOrderId_key, amount_key, discount_key)
}

func (this *Meizu) CheckSign(params ...interface{}) (err error) {
	format := `app_id=%s`
	if buy_amount := this.urlParams.Get("buy_amount"); buy_amount != "" {
		format = format + `&buy_amount=` + buy_amount
	}
	format = format + `&cp_order_id=%s&create_time=%s&notify_id=%s&notify_time=%s&order_id=%s&partner_id=%s&pay_time=%s`

	if pay_type := this.urlParams.Get("pay_type"); pay_type != "" {
		format = format + `&pay_type=` + pay_type
	}
	format = format + `&product_id=%s`

	if product_per_price := this.urlParams.Get("product_per_price"); product_per_price != "" {
		format = format + `&product_per_price=` + product_per_price
	}
	if product_unit := this.urlParams.Get("product_unit"); product_unit != "" {
		format = format + `&product_unit=` + product_unit
	}
	format = format + `&total_price=%s&trade_status=%s&uid=%s`

	if user_info := this.urlParams.Get("user_info"); user_info != "" {
		format = format + `&user_info=` + user_info
	}
	format = format + `:%s`

	this.signContent = fmt.Sprintf(format,
		this.urlParams.Get("app_id"),
		this.urlParams.Get("cp_order_id"),
		this.urlParams.Get("create_time"),
		this.urlParams.Get("notify_id"),
		this.urlParams.Get("notify_time"),
		this.urlParams.Get("order_id"),
		this.urlParams.Get("partner_id"),
		this.urlParams.Get("pay_time"),
		this.urlParams.Get("product_id"),
		this.urlParams.Get("total_price"),
		this.urlParams.Get("trade_status"),
		this.urlParams.Get("uid"),
		this.channelParams["_payKey"])
	this.inputSign = this.urlParams.Get("sign")
	return this.MD5.CheckSign()
}

func (this *Meizu) CheckChannelRet(params ...interface{}) (err error) {
	tradeState := this.urlParams.Get("trade_status") == "3"
	tradeFailDesc := `urlParam(trade_status)!="3"`
	return this.MD5.CheckChannelRet(tradeState, tradeFailDesc)
}

func (this *Meizu) GetResult(params ...interface{}) (ret string) {
	channelRetData := this.channelRetData.(*meizuChannelRet)

	switch this.lastError {
	case err_noerror:
		channelRetData.Code = "200"
		channelRetData.Message = "成功发货"
	case err_checkSign:
		fallthrough
	case err_parseInputParam:
		fallthrough
	case err_payAmountError:
		fallthrough
	case err_tradeFail:
		channelRetData.Code = "120014"
		channelRetData.Message = "发货失败"
	case err_prepareLoginRequest:
		fallthrough
	case err_channelUserIsNotExist:
		fallthrough
	case err_handleOrder:
		channelRetData.Code = "120013"
		channelRetData.Message = "尚未发货"
	default:
		channelRetData.Code = "900000"
		channelRetData.Message = "未知异常"
	}
	jsonByte, _ := json.Marshal(channelRetData)
	ret = string(jsonByte)

	this.MD5.GetResult()
	return
}
