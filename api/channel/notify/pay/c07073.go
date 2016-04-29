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

var C07073UrlKeys []string = []string{"data"}

const (
	err_c07073ParseAppSecret = 14101
)

type C07073 struct {
	Base
	appSecret  string
	c07073Data C07073Data
}

type C07073Data struct {
	Orderid     string `json:"orderid"`
	Gameid      string `json:"gameid"`
	Serverid    string `json:"serverid"`
	Uid         string `json:"uid"`
	Amount      string `json:"amount"`
	Time        uint   `json:"time"`
	Sign        string `json:"sign"`
	ExtendsInfo string `json:"extendsInfo"`
}

func NewC07073(channelId, productId int, urlParams *url.Values) *C07073 {
	ret := new(C07073)
	ret.Init(channelId, productId, urlParams)
	return ret
}

func (this *C07073) Init(channelId, productId int, urlParams *url.Values) {
	this.Base.Init(channelId, productId, urlParams, &C07073UrlKeys)
}

func (this *C07073) parseAppSecret() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_c07073ParseAppSecret
			beego.Trace(err)
		}
	}()
	if this.appSecret, err = this.getPackageParam("C07073_SECRET_KEY"); err != nil {
		return
	}
	return
}

func (this *C07073) parseUrlParam() (err error) {
	defer func() {
		if err != nil {
			beego.Trace(err)
		}
	}()

	data := this.urlParams.Get("data")
	beego.Trace(data)

	if err = json.Unmarshal([]byte(data), &this.c07073Data); err != nil {
		this.callbackRet = err_parseUrlParam
		return
	}

	this.orderId = this.c07073Data.ExtendsInfo
	this.channelUserId = this.c07073Data.Uid
	this.channelOrderId = this.c07073Data.Orderid

	payAmount := 0.0
	if payAmount, err = strconv.ParseFloat(this.c07073Data.Amount, 64); err != nil {
		return err
	} else {
		this.payAmount = int(payAmount * 100)
	}
	return
}

func (this *C07073) ParseParam() (err error) {
	if err = this.parseAppSecret(); err != nil {
		return
	}
	if err = this.parseUrlParam(); err != nil {
		return
	}
	if err = this.Base.ParseParam(); err != nil {
		return
	}
	return
}

func (this *C07073) CheckSign() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_checkSign
			beego.Trace(err)
		}
	}()

	format := "amount=%s&gameid=%s&orderid=%s&serverid=%s&time=%d&uid=%s%s"

	content := fmt.Sprintf(format, this.c07073Data.Amount, this.c07073Data.Gameid,
		this.c07073Data.Orderid, this.c07073Data.Serverid, this.c07073Data.Time,
		this.c07073Data.Uid, this.appSecret)

	urlSign := this.c07073Data.Sign
	if sign := tool.Md5([]byte(content)); sign != urlSign {
		msg := fmt.Sprintf("Sign is invalid, content:%s, sign:%s, urlSign:%s", content, sign, urlSign)
		err = errors.New(msg)
		return
	}
	return
}

func (this *C07073) GetResult() (ret string) {
	beego.Trace(this.callbackRet)
	if this.callbackRet == err_noerror {
		ret = "succ"
	} else {
		ret = "fail"
	}
	return
}
