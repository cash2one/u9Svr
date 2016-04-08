package channelPayNotify

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"net/url"
	"strconv"
	"u9/tool"
)

var mangoUrlKeys []string = []string{"ret_code", "error_msg", "aid", "order_no"}

const (
	err_mangoParsePayKey   = 13401
	err_mangoResultFailure = 13402
)

//芒果玩
type Mango struct {
	Base
	payKey string
}

func NewMango(channelId, productId int, urlParams *url.Values) *Mango {
	ret := new(Mango)
	ret.Init(channelId, productId, urlParams)
	return ret
}

func (this *Mango) Init(channelId, productId int, urlParams *url.Values) {
	this.Base.Init(channelId, productId, urlParams, &mangoUrlKeys)
}

func (this *Mango) parsePayKey() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_mangoParsePayKey
			beego.Trace(err)
		}
	}()
	this.payKey, err = this.getPackageParam("SKEY")
	return
}

func (this *Mango) parseUrlParam() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_parseUrlParam
			beego.Trace(err)
		}
	}()

	this.orderId = this.urlParams.Get("extension_field")
	this.channelUserId = this.urlParams.Get("aid")
	this.channelOrderId = this.urlParams.Get("order_no")

	payAmount := 0.0
	if payAmount, err = strconv.ParseFloat(this.urlParams.Get("money"), 64); err != nil {
		return err
	} else {
		this.payAmount = int(payAmount * 100)
	}
	return
}

func (this *Mango) ParseChannelRet() (err error) {
	if result := this.urlParams.Get("ret_code"); result != "0" {
		this.callbackRet = err_mangoResultFailure
	}
	return
}

func (this *Mango) ParseParam() (err error) {
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

func (this *Mango) CheckSign() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_checkSign
			beego.Trace(err)
		}
	}()

	format := "%s%s%s%s%s%s%s%s%s%s"
	content := fmt.Sprintf(format,
		this.urlParams.Get("ret_code"), this.urlParams.Get("aid"), this.urlParams.Get("gid"), this.urlParams.Get("cid"),
		this.payKey, this.urlParams.Get("ts"), url.QueryEscape(this.urlParams.Get("order_no")),
		url.QueryEscape(this.urlParams.Get("pay_order")), this.urlParams.Get("money"), this.urlParams.Get("pay_type"))

	urlSign := this.urlParams.Get("enc")
	if sign := tool.Md5([]byte(content)); sign != urlSign {
		msg := fmt.Sprintf("Sign is invalid, content:%s, sign:%s, urlSign:%s", content, sign, urlSign)
		err = errors.New(msg)
		return
	}

	return
}

func (this *Mango) GetResult() (ret string) {
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
