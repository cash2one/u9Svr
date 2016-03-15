package channelPayNotify

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"net/url"
	"strconv"
	"u9/tool"
)

var cygameUrlKeys []string = []string{"orderid", "username", "gameid", "paytype",
	"amount", "paytime", "attach", "sign"}

const (
	err_cygameAppKey        = 10101
	err_cygameResultFailure = 10102
)

type CYGame struct {
	Base
	appKey string
}

func NewCYGame(channelId, productId int, urlParams *url.Values) *CYGame {
	ret := new(CYGame)
	ret.Init(channelId, productId, urlParams)
	return ret
}

func (this *CYGame) Init(channelId, productId int, urlParams *url.Values) {
	this.Base.Init(channelId, productId, urlParams, &cygameUrlKeys)
}

func (this *CYGame) parseAppKey() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_cygameAppKey
			beego.Trace(err)
		}
	}()
	this.appKey, err = this.getPackageParam("MG_APPKEY")
	return
}

func (this *CYGame) parseUrlParam() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_parseUrlParam
			beego.Trace(err)
		}
	}()

	this.orderId = this.urlParams.Get("attach")
	this.channelOrderId = this.urlParams.Get("orderid")
	this.channelUserId = this.urlParams.Get("username")
	payAmount := 0.0
	if payAmount, err = strconv.ParseFloat(this.urlParams.Get("amount"), 64); err != nil {
		return err
	} else {
		this.payAmount = int(payAmount * 100)
	}
	return
}

func (this *CYGame) ParseParam() (err error) {
	if err = this.parseUrlParam(); err != nil {
		return
	}
	if err = this.parseAppKey(); err != nil {
		return
	}
	if err = this.Base.ParseParam(); err != nil {
		return
	}
	return
}

func (this *CYGame) CheckSign() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_checkSign
			beego.Trace(err)
		}
	}()

	format := "orderid=%s&username=%s&gameid=%s&roleid=%s&serverid=%s&paytype=%s&amount=%s&paytime=%s&attach=%s&appkey=%s"
	content := fmt.Sprintf(format,
		this.channelOrderId, this.urlParams.Get("username"), this.urlParams.Get("gameid"),
		this.urlParams.Get("roleid"), this.urlParams.Get("serverid"), this.urlParams.Get("paytype"),
		this.urlParams.Get("amount"), this.urlParams.Get("paytime"), this.urlParams.Get("attach"),
		this.appKey)

	urlSign := this.urlParams.Get("sign")
	if sign := tool.Md5([]byte(content)); sign != urlSign {
		msg := fmt.Sprintf("Sign is invalid, content:%s, sign:%s, urlSign%s", content, sign, urlSign)
		err = errors.New(msg)
		return
	}

	return
}

func (this *CYGame) GetResult() (ret string) {
	switch this.callbackRet {
	case err_noerror:
		ret = "success"
	case err_checkSign:
		ret = "errorSign"
		beego.Trace(this.urlParams)
	default:
		ret = "error"
		beego.Trace(this.urlParams)
	}
	return
}
