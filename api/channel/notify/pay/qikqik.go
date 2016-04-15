package channelPayNotify

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"net/url"
	"strconv"
	"u9/tool"
)

var qikqikUrlKeys []string = []string{"uid", "cporder", "money", "order", "cpappid"}

const (
	err_qikqikParsePayKey   = 13501
	err_qikqikResultFailure = 13502
)

//7k7k
type QikQik struct {
	Base
	cpId string
	paySecret string
}

func NewQikQik(channelId, productId int, urlParams *url.Values) *QikQik {
	ret := new(QikQik)
	ret.Init(channelId, productId, urlParams)
	return ret
}

func (this *QikQik) Init(channelId, productId int, urlParams *url.Values) {
	this.Base.Init(channelId, productId, urlParams, &qikqikUrlKeys)
}

func (this *QikQik) parsePayKey() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_qikqikParsePayKey
			beego.Trace(err)
		}
	}()
	this.cpId, err = this.getPackageParam("APPID")
	
	this.paySecret, err = this.getPackageParam("APP_SECRET")
	return
}

func (this *QikQik) parseUrlParam() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_parseUrlParam
			beego.Trace(err)
		}
	}()

	this.orderId = this.urlParams.Get("cporder")
	this.channelUserId = this.urlParams.Get("uid")
	this.channelOrderId = this.urlParams.Get("order")

	payAmount := 0.0
	if payAmount, err = strconv.ParseFloat(this.urlParams.Get("money"), 64); err != nil {
		return err
	} else {
		this.payAmount = int(payAmount * 100)
	}
	return
}

func (this *QikQik) ParseChannelRet() (err error) {
	// if result := this.urlParams.Get("result"); result != "1" {
	// 	this.callbackRet = err_qikqikResultFailure
	// }
	return
}

func (this *QikQik) ParseParam() (err error) {
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

func (this *QikQik) CheckSign() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_checkSign
			beego.Trace(err)
		}
	}()

	format := "%s%s%s%s%s%s"
	content := fmt.Sprintf(format,
		this.channelUserId,this.orderId,this.urlParams.Get("money"),this.channelOrderId,
		this.cpId,this.paySecret)

	urlSign := this.urlParams.Get("sign")
	if sign := tool.Md5([]byte(content)); sign != urlSign {
		msg := fmt.Sprintf("Sign is invalid, content:%s, sign:%s, urlSign:%s", content, sign, urlSign)
		err = errors.New(msg)
		return
	}

	return
}

func (this *QikQik) GetResult() (ret string) {
	if this.callbackRet == err_noerror {
		ret = "success"
	} else {
		ret = "fail"
		beego.Trace(this.callbackRet)
	}
	return
}

