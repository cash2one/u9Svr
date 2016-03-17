package channelPayNotify

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"net/url"
	"strconv"
	"u9/tool"
)

var dangleUrlKeys []string = []string{"result", "money", "order", "mid", "time", "ext", "signature"}

const (
	err_dangleParsePayKey   = 10101
	err_dangleResultFailure = 10102
)

//当乐
type Dangle struct {
	Base
	payKey string
}

func NewDangle(channelId, productId int, urlParams *url.Values) *Dangle {
	ret := new(Dangle)
	ret.Init(channelId, productId, urlParams)
	return ret
}

func (this *Dangle) Init(channelId, productId int, urlParams *url.Values) {
	this.Base.Init(channelId, productId, urlParams, &dangleUrlKeys)
}

func (this *Dangle) parsePayKey() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_dangleParsePayKey
			beego.Trace(err)
		}
	}()
	this.payKey, err = this.getPackageParam("DANGLE_PAYMENT_KEY")
	return
}

func (this *Dangle) parseUrlParam() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_parseUrlParam
			beego.Trace(err)
		}
	}()

	this.orderId = this.urlParams.Get("ext")
	this.channelUserId = this.urlParams.Get("mid")
	this.channelOrderId = this.urlParams.Get("order")

	payAmount := 0.0
	if payAmount, err = strconv.ParseFloat(this.urlParams.Get("money"), 64); err != nil {
		return err
	} else {
		this.payAmount = int(payAmount * 100)
	}
	return
}

func (this *Dangle) ParseChannelRet() (err error) {
	if result := this.urlParams.Get("result"); result != "1" {
		this.callbackRet = err_dangleResultFailure
	}
	return
}

func (this *Dangle) ParseParam() (err error) {
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

func (this *Dangle) CheckSign() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_checkSign
			beego.Trace(err)
		}
	}()

	format := "order=%s&money=%s&mid=%s&time=%s&result=%s&ext=%s&key=%s"
	context := fmt.Sprintf(format,
		this.channelOrderId, this.urlParams.Get("money"),
		this.channelUserId, this.urlParams.Get("time"), this.urlParams.Get("result"),
		this.urlParams.Get("ext"), this.payKey)

	if sign := tool.Md5([]byte(context)); sign != this.urlParams.Get("signature") {
		msg := fmt.Sprintf("Sign is invalid, context:%s, sign:%s", context, sign)
		err = errors.New(msg)
		return
	}
	return
}

func (this *Dangle) GetResult() (ret string) {
	if this.callbackRet == err_noerror {
		ret = "success"
	} else {
		ret = "failure"
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
