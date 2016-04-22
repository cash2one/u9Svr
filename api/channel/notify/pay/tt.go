package channelPayNotify

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"net/url"
	"strconv"
	"u9/tool"
)

var ttUrlKeys []string = []string{"uid", "gameId", "sdkOrderId", "cpOrderId", "payFee", "payResult"}

const (
	err_ttParsePayKey   = 10101
	err_ttResultFailure = 10102
)

//TT语音
type TT struct {
	Base
	payKey string
}

func NewTT(channelId, productId int, urlParams *url.Values) *TT {
	ret := new(TT)
	ret.Init(channelId, productId, urlParams)
	return ret
}

func (this *TT) Init(channelId, productId int, urlParams *url.Values) {
	this.Base.Init(channelId, productId, urlParams, &ttUrlKeys)
}

func (this *TT) parsePayKey() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_ttParsePayKey
			beego.Trace(err)
		}
	}()
	this.payKey, err = this.getPackageParam("TT_PAYMENT_KEY")
	return
}

func (this *TT) parseUrlParam() (err error) {
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

func (this *TT) ParseChannelRet() (err error) {
	if result := this.urlParams.Get("result"); result != "1" {
		this.callbackRet = err_ttResultFailure
	}
	return
}

func (this *TT) ParseParam() (err error) {
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

func (this *TT) CheckSign() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_checkSign
			beego.Trace(err)
		}
	}()

	format := "order=%s&money=%s&mid=%s&time=%s&result=%s&ext=%s&key=%s"
	content := fmt.Sprintf(format,
		this.channelOrderId, this.urlParams.Get("money"),
		this.channelUserId, this.urlParams.Get("time"), this.urlParams.Get("result"),
		this.urlParams.Get("ext"), this.payKey)

	urlSign := this.urlParams.Get("signature")
	if sign := tool.Md5([]byte(content)); sign != urlSign {
		msg := fmt.Sprintf("Sign is invalid, content:%s, sign:%s, urlSign%s", content, sign, urlSign)
		err = errors.New(msg)
		return
	}

	return
}

func (this *TT) GetResult() (ret string) {
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
