package channelPayNotify

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"net/url"
	"strconv"
	"strings"
	"u9/tool"
)

var youlongUrlKeys []string = []string{"orderId", "userName", "amount", "flag"}

const (
	err_youlongParsePayKey = 13301
	// err_youlongResultFailure = 13302
)

//游龙
type YouLong struct {
	Base
	payKey string
}

func NewYouLong(channelId, productId int, urlParams *url.Values) *YouLong {
	ret := new(YouLong)
	ret.Init(channelId, productId, urlParams)
	return ret
}

func (this *YouLong) Init(channelId, productId int, urlParams *url.Values) {
	this.Base.Init(channelId, productId, urlParams, &youlongUrlKeys)
}

func (this *YouLong) parsePayKey() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_youlongParsePayKey
			beego.Trace(err)
		}
	}()
	this.payKey, err = this.getPackageParam("PKEY")
	return
}

func (this *YouLong) parseUrlParam() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_parseUrlParam
			beego.Trace(err)
		}
	}()

	this.orderId = this.urlParams.Get("orderId")
	this.channelUserId = this.urlParams.Get("userName")

	payAmount := 0.0
	if payAmount, err = strconv.ParseFloat(this.urlParams.Get("amount"), 64); err != nil {
		return err
	} else {
		this.payAmount = int(payAmount * 100)
	}
	return
}

func (this *YouLong) ParseChannelRet() (err error) {
	// beego.Trace(this.urlParams)
	// if result := this.urlParams.Get("state"); result != "1" {
	// 	this.callbackRet = err_youlongResultFailure
	// }
	return
}

func (this *YouLong) ParseParam() (err error) {
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

func (this *YouLong) CheckSign() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_checkSign
			beego.Trace(err)
		}
	}()

	format := "%s%s%s%s%s"
	content := fmt.Sprintf(format,
		this.orderId, this.channelUserId,
		this.urlParams.Get("amount"), this.urlParams.Get("extra"),
		this.payKey)

	urlSign := this.urlParams.Get("flag")
	if sign := strings.ToUpper(tool.Md5([]byte(content))); sign != urlSign {
		msg := fmt.Sprintf("Sign is invalid, content:%s, sign:%s, urlSign:%s", content, sign, urlSign)
		err = errors.New(msg)
		return
	}

	return
}

func (this *YouLong) GetResult() (ret string) {
	beego.Trace(this.callbackRet)
	if this.callbackRet == err_noerror {
		ret = "OK"
	} else {
		ret = "NO"
	}
	return
}
