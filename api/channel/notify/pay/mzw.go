package channelPayNotify

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"net/url"
	"strconv"
	"u9/tool"
)

var mzwUrlKeys []string = []string{"appkey", "orderID", "money", "uid", "extern", "sign"}

const (
	err_mzwParsePayKey   = 11301
)

//拇指玩
type MZW struct {
	Base
	appkey string
	payKey string
}

func NewMZW(channelId, productId int, urlParams *url.Values) *MZW {
	ret := new(MZW)
	ret.Init(channelId, productId, urlParams)
	return ret
}

func (this *MZW) Init(channelId, productId int, urlParams *url.Values) {
	this.Base.Init(channelId, productId, urlParams, &mzwUrlKeys)
}

func (this *MZW) parsePayKey() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_mzwParsePayKey
			beego.Trace(err)
		}
	}()
	this.appkey, err = this.getPackageParam("MZWAPPKEY")
	this.payKey, err = this.getPackageParam("MZWPAYKEY")
	return
}

func (this *MZW) parseUrlParam() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_parseUrlParam
			beego.Trace(err)
		}
	}()

	this.orderId = this.urlParams.Get("extern")
	this.channelUserId = this.urlParams.Get("uid")
	this.channelOrderId = this.urlParams.Get("orderID")

	payAmount := 0.0
	if payAmount, err = strconv.ParseFloat(this.urlParams.Get("money"), 64); err != nil {
		return err
	} else {
		this.payAmount = int(payAmount * 100)
	}
	return
}

func (this *MZW) ParseChannelRet() (err error) {
	// if result := this.urlParams.Get("result"); result != "1" {
	// 	this.callbackRet = err_mzwResultFailure
	// }
	return
}

func (this *MZW) ParseParam() (err error) {
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

func (this *MZW) CheckSign() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_checkSign
			beego.Trace(err)
		}
	}()

	format := "%s%s%s%s%s%s%s%s%s"
	content := fmt.Sprintf(format,
		this.appkey, this.channelOrderId,this.urlParams.Get("productName"),this.urlParams.Get("productDesc"),
		this.urlParams.Get("productID"), this.urlParams.Get("money"),this.channelUserId,this.urlParams.Get("extern"), 
		this.payKey)

	urlSign := this.urlParams.Get("sign")
	if sign := tool.Md5([]byte(content)); sign != urlSign {
		msg := fmt.Sprintf("Sign is invalid, content:%s, sign:%s, urlSign%s", content, sign, urlSign)
		err = errors.New(msg)
		return
	}

	return
}

func (this *MZW) GetResult() (ret string) {
	if this.callbackRet == err_noerror {
		ret = "SUCCESS"
	} else {
		ret = "FAILURE"
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
