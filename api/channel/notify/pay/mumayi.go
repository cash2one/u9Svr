package channelPayNotify

import (
	// "errors"
	"fmt"
	"github.com/astaxie/beego"
	"net/url"
	"strconv"
	"u9/tool"
	"strings"
)

var mumayiUrlKeys []string = []string{"uid", "orderID", "productPrice", "orderTime", "tradeSign", "tradeState"}

const (
	err_mumayiParsePayKey   = 10101
	err_mumayiResultFailure = 10102
)

//木蚂蚁
type MuMaYi struct {
	Base
	payKey string
}

func NewMuMaYi(channelId, productId int, urlParams *url.Values) *MuMaYi {
	ret := new(MuMaYi)
	ret.Init(channelId, productId, urlParams)
	return ret
}

func (this *MuMaYi) Init(channelId, productId int, urlParams *url.Values) {
	this.Base.Init(channelId, productId, urlParams, &mumayiUrlKeys)
}

func (this *MuMaYi) parsePayKey() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_mumayiParsePayKey
			beego.Trace(err)
		}
	}()
	this.payKey, err = this.getPackageParam("MUMAYI_APPKEY")
	return
}

func (this *MuMaYi) parseUrlParam() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_parseUrlParam
			beego.Trace(err)
		}
	}()

	this.orderId = this.urlParams.Get("productDesc")
	this.channelUserId = this.urlParams.Get("uid")
	this.channelOrderId = this.urlParams.Get("orderID")

	payAmount := 0.0
	if payAmount, err = strconv.ParseFloat(this.urlParams.Get("productPrice"), 64); err != nil {
		return err
	} else {
		this.payAmount = int(payAmount * 100)
	}
	return
}

func (this *MuMaYi) ParseChannelRet() (err error) {
	if result := this.urlParams.Get("tradeState"); result != "success" {
		this.callbackRet = err_mumayiResultFailure
	}
	return
}

func (this *MuMaYi) ParseParam() (err error) {
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

func (this *MuMaYi) CheckSign() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_checkSign
			beego.Trace(err)
		}
	}()

	var result string
	urlSign := this.urlParams.Get("tradeSign")
	if result,err = tool.MMYSign(urlSign,this.payKey,this.channelOrderId);err != nil{
		msg := fmt.Sprintf("mumayi Sign is erro,  sign:%s, erro:%s", urlSign,err )
		beego.Error(msg)
	}
	result = strings.TrimSpace(result)
	if result != "true"{
		this.callbackRet = err_checkSign
		beego.Trace("mmy check:", result, "urlSign:", urlSign, "payKey:",this.payKey, "channelOrderId:", this.channelOrderId)
	}else{
		beego.Trace("mmy check:", result)
	}
	// if sign := tool.Md5([]byte(content)); sign != urlSign {
	// 	msg := fmt.Sprintf("Sign is invalid, content:%s, sign:%s, urlSign%s", content, sign, urlSign)
	// 	err = errors.New(msg)
	// 	return
	// }

	return
}

func (this *MuMaYi) GetResult() (ret string) {
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
