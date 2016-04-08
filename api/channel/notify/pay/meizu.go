package channelPayNotify

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"net/url"
	"strconv"
	"u9/tool"
)

var meizuUrlKeys []string = []string{"notify_time", "notify_id", "order_id", "app_id", "uid",
	"partner_id", "cp_order_id", "product_id", "total_price", "trade_status", "create_time",
	"pay_time", "sign", "sign_type"}

const (
	err_meizuParsePayKey   = 11201
	err_meizuResultFailure = 11202
)

type Meizu struct {
	Base
	appSecret string
}

func NewMeizu(channelId, productId int, urlParams *url.Values) *Meizu {
	ret := new(Meizu)
	ret.Init(channelId, productId, urlParams)
	return ret
}

func (this *Meizu) Init(channelId, productId int, urlParams *url.Values) {
	this.Base.Init(channelId, productId, urlParams, &meizuUrlKeys)
}

func (this *Meizu) parsePayKey() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_meizuParsePayKey
			beego.Trace(err)
		}
	}()
	this.appSecret, err = this.getPackageParam("MEIZU_APPSECRET")
	return
}

func (this *Meizu) parseUrlParam() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_parseUrlParam
			beego.Trace(err)
		}
	}()

	this.orderId = this.urlParams.Get("cp_order_id")
	this.channelUserId = this.urlParams.Get("uid")
	this.channelOrderId = this.urlParams.Get("order_id")

	payAmount := 0.0
	if payAmount, err = strconv.ParseFloat(this.urlParams.Get("total_price"), 64); err != nil {
		return err
	} else {
		this.payAmount = int(payAmount * 100)
	}
	return
}

func (this *Meizu) ParseChannelRet() (err error) {
	if result := this.urlParams.Get("trade_status"); result != "3" {
		beego.Trace(result)
		this.callbackRet = err_meizuResultFailure
	}
	return
}

func (this *Meizu) ParseParam() (err error) {
	if err = this.parseUrlParam(); err != nil {
		return
	}
	if err = this.parsePayKey(); err != nil {
		return
	}
	if err = this.Base.ParseParam(); err != nil {
		return
	}
	return
}

func (this *Meizu) CheckSign() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_checkSign
			beego.Trace(err)
		}
	}()

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

	context := fmt.Sprintf(format, this.urlParams.Get("app_id"), this.urlParams.Get("cp_order_id"),
		this.urlParams.Get("create_time"), this.urlParams.Get("notify_id"), this.urlParams.Get("notify_time"),
		this.urlParams.Get("order_id"), this.urlParams.Get("partner_id"), this.urlParams.Get("pay_time"),
		this.urlParams.Get("product_id"), this.urlParams.Get("total_price"), this.urlParams.Get("trade_status"),
		this.urlParams.Get("uid"), this.appSecret)

	if sign := tool.Md5([]byte(context)); sign != this.urlParams.Get("sign") {
		msg := fmt.Sprintf("Sign is invalid, context:%s, sign:%s", context, sign)
		err = errors.New(msg)
		return
	}
	return
}

func (this *Meizu) GetResult() (ret string) {
	type meizuRet struct {
		Code     string `json:"code"` //200 成功发货；120013 尚未发货；120014 发货失败; 900000 未知异常
		Message  string `json:"message"`
		Value    string `json:"value"`
		Redirect string `json:"redirect"`
	}
	jsonRet := meizuRet{
		Code:     "200",
		Message:  "",
		Value:    "",
		Redirect: "",
	}

	switch this.callbackRet {
	case err_noerror:
		jsonRet.Code = "200"
		jsonRet.Message = "成功发货"
	case err_checkSign:
		fallthrough
	case err_parseUrlParam:
		fallthrough
	case err_payAmountError:
		fallthrough
	case err_meizuResultFailure:
		jsonRet.Code = "120014"
		jsonRet.Message = "发货失败"
	case err_parseLoginRequest:
		fallthrough
	case err_channelUserIsNotExist:
		fallthrough
	case err_handleOrder:
		jsonRet.Code = "120013"
		jsonRet.Message = "尚未发货"
	default:
		jsonRet.Code = "900000"
		jsonRet.Message = "未知异常"
	}
	data, _ := json.Marshal(jsonRet)
	ret = string(data)
	return
}
