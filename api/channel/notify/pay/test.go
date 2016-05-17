package channelPayNotify

import (
	"github.com/astaxie/beego"
	"net/url"
)

type Test struct {
	Base
}

func NewTest(channelId, productId int, urlParams *url.Values) *Test {
	ret := new(Test)
	ret.Init(channelId, productId, urlParams)
	return ret
}

func (this *Test) Init(channelId, productId int, urlParams *url.Values) {
	this.Base.Init(channelId, productId, urlParams, &emptyUrlKeys)
}

func (this *Test) ParseChannelRet() (err error) {
	return
}

func (this *Test) parseUrlParam() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_parseUrlParam
			beego.Trace(err)
		}
	}()
	return
}

func (this *Test) ParseParam() (err error) {
	if err = this.parseUrlParam(); err != nil {
		return
	}
	return
}

func (this *Test) CheckSign() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_checkSign
			beego.Trace(err)
		}
	}()
	return
}

func (this *Test) GetResult() (ret string) {
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
