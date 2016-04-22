package channelPayNotify

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"net/url"
	"strconv"
	"u9/tool"
)

var letvUrlKeys []string = []string{"app_id", "lepay_order_no", "letv_user_id", "out_trade_no",
	"pay_time", "price", "product_id", "sign", "sign_type", "trade_result", "version",
	"cooperator_order_no"}

const (
	err_letvParsePayKey   = 12601
	err_letvResultFailure = 12602
)

//乐视
type LeTV struct {
	Base
	appid  string
	payKey string
}

func NewLeTV(channelId, productId int, urlParams *url.Values) *LeTV {
	ret := new(LeTV)
	ret.Init(channelId, productId, urlParams)
	return ret
}

func (this *LeTV) Init(channelId, productId int, urlParams *url.Values) {
	this.Base.Init(channelId, productId, urlParams, &letvUrlKeys)
}

func (this *LeTV) parsePayKey() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_letvParsePayKey
			beego.Trace(err)
		}
	}()
	this.payKey, err = this.getPackageParam("lepay_secretkey")
	this.appid, err = this.getPackageParam("lepay_appid")
	return
}

func (this *LeTV) parseUrlParam() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_parseUrlParam
			beego.Trace(err)
		}
	}()

	this.orderId = this.urlParams.Get("cooperator_order_no")
	this.channelUserId = this.urlParams.Get("letv_user_id")
	this.channelOrderId = this.urlParams.Get("out_trade_no")

	payAmount := 0.0
	if payAmount, err = strconv.ParseFloat(this.urlParams.Get("price"), 64); err != nil {
		return err
	} else {
		this.payAmount = int(payAmount * 100)
	}
	return
}

func (this *LeTV) ParseChannelRet() (err error) {
	if result := this.urlParams.Get("trade_result"); result != "TRADE_SUCCESS" {
		this.callbackRet = err_letvResultFailure
	}
	return
}

func (this *LeTV) ParseParam() (err error) {
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

func (this *LeTV) CheckSign() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_checkSign
			beego.Trace(err)
		}
	}()
 switch this.productId{
 	//破阵无双SDK 老版本  其他均为新版本
 	case 1001:
		format := "app_id=%s&lepay_order_no=%s&letv_user_id=%s&out_trade_no=%s&pay_time=%s&price=%s&product_id=%s&sign_type=%s&trade_result=%s&version=%s&key=%s"
		context := fmt.Sprintf(format,
			this.appid, this.urlParams.Get("lepay_order_no"),
			this.channelUserId, this.channelOrderId, this.urlParams.Get("pay_time"),
			this.urlParams.Get("price"), this.urlParams.Get("product_id"), "MD5",
			this.urlParams.Get("trade_result"), this.urlParams.Get("version"), this.payKey)

		if sign := tool.Md5([]byte(context)); sign != this.urlParams.Get("sign") {
			msg := fmt.Sprintf("Sign is invalid, context:%s, sign:%s", context, sign)
			err = errors.New(msg)
			return
		}
		return

	default :
		format := "app_id=%s&lepay_order_no=%s&letv_user_id=%s&out_trade_no=%s&pay_time=%s&price=%s&product_id=%s&sign_type=%s&trade_result=%s&version=%s&cooperator_order_no=%s&extra_info=%s&original_price=%s&key=%s"
		context := fmt.Sprintf(format,
			this.appid, this.urlParams.Get("lepay_order_no"),
			this.channelUserId, this.channelOrderId, this.urlParams.Get("pay_time"),
			this.urlParams.Get("price"), this.urlParams.Get("product_id"), "MD5",
			this.urlParams.Get("trade_result"), this.urlParams.Get("version"),this.orderId,
			this.urlParams.Get("extra_info"),this.urlParams.Get("original_price"),this.payKey)

		if sign := tool.Md5([]byte(context)); sign != this.urlParams.Get("sign") {
			msg := fmt.Sprintf("Sign is invalid, context:%s, sign:%s", context, sign)
			err = errors.New(msg)
			return
		}
		return
	}

}

func (this *LeTV) GetResult() (ret string) {
	if this.callbackRet == err_noerror {
		ret = "success"
	} else {
		ret = "fail"
	}
	return
}
