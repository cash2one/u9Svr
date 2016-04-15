package channelPayNotify

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"net/url"
	"strconv"
	"u9/tool"
)

var sogouUrlKeys []string = []string{"gid", "sid", "uid", "oid", "date", "amount1","realAmount","auth"}

const (
	err_sogouParsePayKey   = 11501
	err_sogouResultFailure = 11502
	err_sogouOk = "OK"
	err_sogouParam = "ERR_100"
	err_sogouSign = "ERR_200"
	// err_sogouUid = "ERR_300"
	// err_sogouIP = "ERR_400"
	// err_sogouOther = "ERR_500"
)

//当乐
type Sogou struct {
	Base
	sogouPayResult string
	gid string
	payKey string
}

func NewSogou(channelId, productId int, urlParams *url.Values) *Sogou {
	ret := new(Sogou)
	ret.Init(channelId, productId, urlParams)
	return ret
}

func (this *Sogou) Init(channelId, productId int, urlParams *url.Values) {
	this.Base.Init(channelId, productId, urlParams, &sogouUrlKeys)
}

func (this *Sogou) parsePayKey() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_sogouParsePayKey
			beego.Trace(err)
		}
	}()

	this.gid, err = this.getPackageParam("SOGOU_GAMEID")
	this.payKey, err = this.getPackageParam("SOGOU_PAYKEY")
	return
}

func (this *Sogou) parseUrlParam() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_parseUrlParam
			beego.Trace(err)
		}
	}()

	this.orderId = this.urlParams.Get("appdata")
	this.channelUserId = this.urlParams.Get("uid")
	this.channelOrderId = this.urlParams.Get("oid")

	payAmount := 0.0
	if payAmount, err = strconv.ParseFloat(this.urlParams.Get("realAmount"), 64); err != nil {
		return err
	} else {
		this.payAmount = int(payAmount * 100)
	}
	return
}

func (this *Sogou) ParseChannelRet() (err error) {
	// if result := this.urlParams.Get("result"); result != "1" {
	// 	this.callbackRet = err_sogouResultFailure
	// }
	return
}

func (this *Sogou) ParseParam() (err error) {
	if err = this.parseUrlParam(); err != nil {
		return
	}
	if err = this.parsePayKey(); err != nil {
		return
	}
	if err = this.Base.ParseParam(); err != nil {
		this.sogouPayResult = err_sogouParam
		return
	}
	return
}

func (this *Sogou) CheckSign() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_checkSign
			this.sogouPayResult = err_sogouSign
			beego.Trace(err)
		}
	}()

	format := "amount1=%s&amount2=%s&appdata=%s&date=%s&gid=%s&oid=%s&realAmount=%s&role=%s&sid=%s&time=%s&uid=%s&%s"
	content := fmt.Sprintf(format,
		this.urlParams.Get("amount1"), this.urlParams.Get("amount2"),this.orderId,this.urlParams.Get("date"),
		this.gid,this.channelOrderId, this.urlParams.Get("realAmount"), this.urlParams.Get("role"),
		this.urlParams.Get("sid"),this.urlParams.Get("time"), this.channelUserId,this.payKey)

	urlSign := this.urlParams.Get("auth")
	if sign := tool.Md5([]byte(content)); sign != urlSign {
		msg := fmt.Sprintf("Sign is invalid, content:%s, sign:%s, urlSign:%s", content, sign, urlSign)
		err = errors.New(msg)
		return
	}

	return
}

func (this *Sogou) GetResult() (ret string) {
	if this.callbackRet == err_noerror {
		ret = "OK"
	} else {
		ret = this.sogouPayResult
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
