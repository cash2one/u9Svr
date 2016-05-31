package channelPayNotify

import (
	"github.com/astaxie/beego"
	"net/url"
	"strconv"
)

var testUrlKeys []string = []string{"money", "order", "mid",
	"time", "ext", "signature", "result"}

type Test struct {
	Base
}

func NewTest(channelId, productId int, urlParams *url.Values) *Test {
	ret := new(Test)
	ret.Init(channelId, productId, urlParams)
	return ret
}

func (this *Test) Init(channelId, productId int, urlParams *url.Values) {
	this.Base.Init(channelId, productId, urlParams, &testUrlKeys)
}

func (this *Test) parseUrlParam() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_parseUrlParam
			beego.Error(err)
			beego.Error(this.urlParams)
		}
	}()

	this.orderId = this.urlParams.Get("ext")
	this.channelOrderId = this.urlParams.Get("order")
	this.channelUserId = this.urlParams.Get("mid")

	payAmount := 0
	if payAmount, err = strconv.Atoi(this.urlParams.Get("money")); err != nil {
		return err
	} else {
		this.payAmount = int(payAmount * 100)
	}
	return
}

func (this *Test) ParseParam() (err error) {
	if err = this.parseUrlParam(); err != nil {
		return
	}
	if err = this.Base.ParseParam(); err != nil {
		return
	}
	return
}

func (this *Test) ParseChannelRet() (err error) {
	if err = this.Base.ParseChannelRet(); err != nil {
		return
	}
	if this.urlParams.Get("result") != "0" {
		this.callbackRet = err_callbackFail
		return
	}
	return
}

func (this *Test) CheckSign() (err error) {
	return
}

func (this *Test) GetResult() (ret string) {
	beego.Trace("callbackRet:" + strconv.Itoa(this.callbackRet))
	if this.callbackRet == err_noerror {
		ret = "{\"code\":\"0\"}"
	} else {
		ret = "{\"code\":\"1\"}"
	}
	return
}

/*
  test url:
  http://192.168.0.185/api/channelPayNotify/1000/101/?result=1
*/
