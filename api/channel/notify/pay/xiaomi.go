package channelPayNotify

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"net/url"
	"strconv"
	"u9/models"
	"u9/tool"
)

var xiaomiUrlKeys []string = []string{"appId", "cpOrderId", "uid", "orderId", "orderStatus", "payFee",
	"productCode", "productName", "productCount", "payTime", "signature"}

const (
	err_xiaomiParseSecretKey  = 11801
	err_xiaomiResultFailure   = 11802
	err_xiaomiAppidIsNotExist = 11803
	err_xiaomiUserIdIsError   = 11804
)

type Xiaomi struct {
	Base
	secretKey string
}

func NewXiaomi(channelId, productId int, urlParams *url.Values) *Xiaomi {
	ret := new(Xiaomi)
	ret.Init(channelId, productId, urlParams)
	return ret
}

func (this *Xiaomi) Init(channelId, productId int, urlParams *url.Values) {
	this.Base.Init(channelId, productId, urlParams, &xiaomiUrlKeys)
}

func (this *Xiaomi) parseSecretKey() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_xiaomiParseSecretKey
			beego.Trace(err)
		}
	}()
	if this.secretKey, err = this.getPackageParam("XIAOMI_SECRETKEY"); err != nil {
		return
	}
	return
}

func (this *Xiaomi) checkUserId() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_xiaomiUserIdIsError
			beego.Trace(err)
		}
	}()

	userId := models.GenerateUserId(this.channelId, this.productId, this.channelUserId)
	beego.Trace(userId, ":", len(userId))
	beego.Trace(this.orderRequest.UserId, ":", len(this.orderRequest.UserId))
	if userId != this.orderRequest.UserId {
		format := `orderRequest's userId(%s) is match url params(channelId(%d) productId(%d) channelUserId(%s))`
		err = errors.New(fmt.Sprintf(format, userId, this.channelId, this.productId, this.channelUserId))
		return
	}
	return
}

func (this *Xiaomi) parseUrlParam() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_parseUrlParam
			beego.Trace(err)
		}
	}()

	this.orderId = this.urlParams.Get("cpOrderId")
	this.channelUserId = this.urlParams.Get("uid")
	this.channelOrderId = this.urlParams.Get("orderId")

	payAmount := 0
	if payAmount, err = strconv.Atoi(this.urlParams.Get("payFee")); err != nil {
		return
	} else {
		this.payAmount = payAmount
	}
	return
}

func (this *Xiaomi) ParseChannelRet() (err error) {
	if this.urlParams.Get("orderStatus") != "TRADE_SUCCESS" {
		err = errors.New("orderStatus isn't equal TRADE_SUCCESS")
		this.callbackRet = err_xiaomiResultFailure
		return
	}

	//appid error 1515
	appId := ""
	if appId, err = this.getPackageParam("XIAOMI_APPID"); err != nil {
		this.callbackRet = err_xiaomiAppidIsNotExist
		return
	}
	if appId != this.urlParams.Get("appId") {
		format := "packageParam appid(%s) is equal urlParam appid(%s)"
		err = errors.New(fmt.Sprintf(format, appId, this.urlParams.Get("appId")))
		this.callbackRet = err_xiaomiAppidIsNotExist
		return
	}
	return
}

func (this *Xiaomi) ParseParam() (err error) {
	if err = this.parseUrlParam(); err != nil {
		return
	}
	if err = this.Base.ParseParam(); err != nil {
		return
	}
	if err = this.checkUserId(); err != nil {
		return
	}
	if err = this.parseSecretKey(); err != nil {
		return
	}
	return
}

func (this *Xiaomi) CheckSign() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_checkSign
			beego.Trace(this.urlParams)
			beego.Trace(err)
		}
	}()

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

	appId := this.urlParams.Get("appId")
	orderStatus := this.urlParams.Get("orderStatus")
	payFee := this.urlParams.Get("payFee")
	payTime := this.urlParams.Get("payTime")
	productCode := this.urlParams.Get("productCode")
	productCount := this.urlParams.Get("productCount")
	productName := this.urlParams.Get("productName")

	context := fmt.Sprintf(format, appId, this.orderId, this.channelOrderId, orderStatus,
		payFee, payTime, productCode, productCount, productName, this.channelUserId)

	sign := fmt.Sprintf("%x", string(tool.HmacSHA1Encrypt(context, this.secretKey)))

	urlSign := this.urlParams.Get("signature")
	if sign != urlSign {
		msg := fmt.Sprintf("Sign is invalid, context:%s, sign:%s, urlSign:%s", context, sign, urlSign)
		err = errors.New(msg)
		return
	}

	return
}

func (this *Xiaomi) GetResult() (ret string) {
	type XiaomiRet struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errMsg"`
	}

	xiaomiRet := new(XiaomiRet)

	switch this.callbackRet {
	case err_noerror:
		xiaomiRet.ErrCode = 200
		xiaomiRet.ErrMsg = "success"
	case err_orderIsNotExist:
		xiaomiRet.ErrCode = 1506
		xiaomiRet.ErrMsg = "cpOrderId error"
	case 303:
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
		xiaomiRet.ErrCode = this.callbackRet
		xiaomiRet.ErrMsg = "other error"
	}
	data, _ := json.Marshal(xiaomiRet)
	ret = string(data)
	return
}
