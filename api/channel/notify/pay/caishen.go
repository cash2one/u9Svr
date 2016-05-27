package channelPayNotify

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"net/url"
	"strconv"
	"u9/tool"
)

var caishenUrlKeys []string = []string{"product_id", "order_id", "price", "game_uid", "u_id", "xcs_order", "game_id", "sign"}

const (
	err_caishenParsePayKey     = 13601
	err_caishenGetUrlParam     = "0001"
	success_caishenGetUrlParam = "0000"
	err_caishenSign            = "0002"
	// err_caishenResultFailure = 13602
)

//小财神
type Xcs struct {
	Base
	payResult string
	payKey    string
}

func NewXcs(channelId, productId int, urlParams *url.Values) *Xcs {
	ret := new(Xcs)
	ret.Init(channelId, productId, urlParams)
	return ret
}

func (this *Xcs) Init(channelId, productId int, urlParams *url.Values) {
	this.Base.Init(channelId, productId, urlParams, &caishenUrlKeys)
}

func (this *Xcs) parsePayKey() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_caishenParsePayKey
			beego.Trace(err)
		}
	}()
	this.payKey, err = this.getPackageParam("PAYKEY")
	return
}

func (this *Xcs) parseUrlParam() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_parseUrlParam
			beego.Trace(err)
		}
	}()

	this.orderId = this.urlParams.Get("order_id")
	this.channelUserId = this.urlParams.Get("u_id")
	this.channelOrderId = this.urlParams.Get("xcs_order")

	payAmount := 0.0
	if payAmount, err = strconv.ParseFloat(this.urlParams.Get("price"), 64); err != nil {
		return err
	} else {
		this.payAmount = int(payAmount * 100)
	}
	return
}

func (this *Xcs) ParseChannelRet() (err error) {
	// if result := this.urlParams.Get("result"); result != "1" {
	// 	this.callbackRet = err_caishenResultFailure
	// }
	return
}

func (this *Xcs) ParseParam() (err error) {
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

func (this *Xcs) CheckSign() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_checkSign
			beego.Trace(err)
		}
	}()

	format := "%s_%s_%s_%s_%s_%s_%s_%s"
	content := fmt.Sprintf(format, this.orderId, this.urlParams.Get("product_id"),
		this.urlParams.Get("price"), this.urlParams.Get("game_uid"), this.channelUserId,
		this.channelOrderId, this.urlParams.Get("game_id"), this.payKey)

	urlSign := this.urlParams.Get("sign")
	if sign := tool.Md5([]byte(content)); sign != urlSign {
		msg := fmt.Sprintf("Sign is invalid, content:%s, sign:%s, urlSign:%s", content, sign, urlSign)
		this.payResult = err_caishenSign
		err = errors.New(msg)
		return
	}

	return
}

func (this *Xcs) GetResult() (ret string) {

	if this.callbackRet == err_noerror {
		ret = success_caishenGetUrlParam
	} else {
		ret = this.payResult
		beego.Trace(this.callbackRet)
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
