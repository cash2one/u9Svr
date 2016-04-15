package channelPayNotify

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"net/url"
	"strconv"
	"u9/tool"
	"u9/models"
)

var pptvUrlKeys []string = []string{"username", "oid", "amount", "extra", "time", "sign"}

const (
	err_pptvParsePayKey   = 13701
	err_pptvResultFailure = 13702
	err_pptvCheckUserId = 13703
	success_pptvPay = `{"code":"1","message":"success"}`
	err_pptvPayMsg = `{"code":"2","message":"pay_fail"}`
	err_pptvOrderMsg = `{"code":"3","message":"order_fail"}`
	err_pptvSignMsg = `{"code":"4","message":"sign_fail"}`
	err_pptvUserIdMsg = `{"code":"5","message":"user_fail"}`
	err_pptvAmountMsg = `{"code":"6","message":"money_fail"}`
)

//PPTV
type PPTV struct {
	Base
	payKey string
}

func NewPPTV(channelId, productId int, urlParams *url.Values) *PPTV {
	ret := new(PPTV)
	ret.Init(channelId, productId, urlParams)
	return ret
}

func (this *PPTV) Init(channelId, productId int, urlParams *url.Values) {
	this.Base.Init(channelId, productId, urlParams, &pptvUrlKeys)
}

func (this *PPTV) parsePayKey() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_pptvParsePayKey
			beego.Trace(err)
		}
	}()
	this.payKey, err = this.getPackageParam("PAYKEY")
	return
}

func (this *PPTV) parseUrlParam() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_parseUrlParam
			beego.Trace(err)
		}
	}()

	this.orderId = this.urlParams.Get("extra")
	this.channelUserId = this.urlParams.Get("username")
	this.channelOrderId = this.urlParams.Get("oid")

	payAmount := 0.0
	if payAmount, err = strconv.ParseFloat(this.urlParams.Get("amount"), 64); err != nil {
		return err
	} else {
		this.payAmount = int(payAmount * 100)
	}
	return
}

func (this *PPTV) ParseChannelRet() (err error) {
	// if result := this.urlParams.Get("code"); result != "1" {
	// 	this.callbackRet = err_pptvResultFailure
	// }
	return
}

func (this *PPTV) checkUserId() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_pptvCheckUserId
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

func (this *PPTV) ParseParam() (err error) {
	if err = this.parseUrlParam(); err != nil {
		return
	}
	if err = this.parsePayKey(); err != nil {
		return
	}
	if err = this.Base.ParseParam(); err != nil {
		return
	}
	if err = this.checkUserId();err !=nil{
		return 
	}
	return
}

func (this *PPTV) CheckSign() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_checkSign
			beego.Trace(err)
		}
	}()

	format := "%s%s%s%s%s%s%s"
	content := fmt.Sprintf(format,
		this.urlParams.Get("sid"),this.channelUserId,this.urlParams.Get("roid"),
		this.channelOrderId,this.urlParams.Get("amount"),this.urlParams.Get("time"),this.payKey)

	urlSign := this.urlParams.Get("sign")
	if sign := tool.Md5([]byte(content)); sign != urlSign {
		msg := fmt.Sprintf("Sign is invalid, content:%s, sign:%s, urlSign%s", content, sign, urlSign)
		err = errors.New(msg)
		return
	}

	return
}

func (this *PPTV) GetResult() (ret string) {
	beego.Trace(this.callbackRet)
	if this.callbackRet == err_noerror {
		ret = success_pptvPay
	} else {
		switch this.callbackRet{
		case err_checkSign:
		ret = err_pptvSignMsg
		case err_pptvCheckUserId:
		ret = err_pptvUserIdMsg
		default :
		ret = err_pptvPayMsg
		}
	}
	return
}

/*
  signature rule: md5("order=xxxx&money=xxxx&mid=xxxx&time=xxxx&result=x&ext=xxx&key=xxxx")
  test url:
  http://192.168.0.185/api/channelPayNotify/1000/101/?
  order=test20160116172500359&
  money=100.00&
  mid=test10086001&
  time=20160116172500&
  result=1&
  ext=game20160116175128772&
  signature=8f00a109716e819bfe0afb695c1addf1
*/
