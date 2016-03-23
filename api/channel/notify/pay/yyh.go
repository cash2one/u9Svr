package channelPayNotify

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"net/url"
	"strconv"
	"u9/tool"
)

var yyhUrlKeys []string = []string{"transdata", "sign"}

const (
	err_yyhParsePayKey   = 10101
	err_yyhResultFailure = 10102
)

//应用汇
type YYH struct {
	Base
	payKey    string
	transData TransData
}
type TransData struct {
	Exorderno string `json:"exorderno"`
	Transid   string `json:"transid"`
	Appid     string `json:"appid"`
	Waresid   string `json:"waresid"`
	Feetype   string `json:"feetype"`
	Money     string `json:"money"`
	Count     string `json:"count"`
	Result    string `json:"result"`
	Transtype string `json:"transtype"`
	Transtime string `json:"transtime"`
	Cpprivate string `json:"cpprivate"`
}

func NewYYH(channelId, productId int, urlParams *url.Values) *YYH {
	ret := new(YYH)
	ret.Init(channelId, productId, urlParams)
	return ret
}

func (this *YYH) Init(channelId, productId int, urlParams *url.Values) {
	this.Base.Init(channelId, productId, urlParams, &yyhUrlKeys)
}

func (this *YYH) parsePayKey() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_yyhParsePayKey
			beego.Trace(err)
		}
	}()
	this.payKey, err = this.getPackageParam("YYH_PAYMENT_KEY")
	return
}

func (this *YYH) parseUrlParam() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_parseUrlParam
			beego.Trace(err)
		}
	}()
	json.Unmarshal(this.urlParams.Get(transdata), &this.transData)
	this.orderId = this.transData.Exorderno
	this.channelUserId = this.urlParams.Get("mid")
	this.channelOrderId = this.transData.Transid

	payAmount := 0.0
	if payAmount, err = strconv.ParseFloat(this.transData.Money, 64); err != nil {
		return err
	} else {
		this.payAmount = int(payAmount)
	}
	return
}

func (this *YYH) ParseChannelRet() (err error) {
	if result := this.urlParams.Get("result"); result != "1" {
		this.callbackRet = err_yyhResultFailure
	}
	return
}

func (this *YYH) ParseParam() (err error) {
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

func (this *YYH) CheckSign() (err error) {
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

func (this *YYH) GetResult() (ret string) {
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
