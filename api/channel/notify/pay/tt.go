package channelPayNotify

import (
	"bytes"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"io"
	"net/url"
	"strconv"
	"u9/tool"
)

var ttUrlKeys []string = []string{}

const (
	err_ttParsePayKey   = 14001
	err_ttResultFailure = 14002
	err_ttParseBody     = 14003
)

//tt
type TT struct {
	Base
	payKey         string
	tt_Sign        string
	tt_result      TT_Result
	tt_contentBody string
	ctx            *context.Context
}

type TT_Result struct {
	Uid        int    `json:"uid"`
	GameId     int    `json:"gameId"`
	SDKOrderId string `json:"sdkOrderId"`
	CpOrderId  string `json:"cpOrderId"`
	PayFee     string `json:"payFee"`
	PayResult  string `json:"payResult"`
	PayDate    string `json:"payDate"`
	ExInfo     string `json:"exInfo"`
}

var (
	ttRsaPublicKey *rsa.PublicKey
)

func NewTT(channelId, productId int, urlParams *url.Values, ctx *context.Context) *TT {
	ret := new(TT)
	ret.Init(channelId, productId, urlParams, ctx)
	return ret
}

func (this *TT) Init(channelId, productId int, urlParams *url.Values, ctx *context.Context) {
	this.Base.Init(channelId, productId, urlParams, &ttUrlKeys)
	this.ctx = ctx
}

func (this *TT) parsePayKey() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_ttParsePayKey
			beego.Error(err)
		}
	}()
	this.payKey, err = this.getPackageParam("TT_SDK_PAYKEY")
	return
}

func (this *TT) parseUrlParam() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_parseUrlParam
			beego.Error(err)
		}
	}()

	beego.Trace(this.tt_result)
	this.orderId = this.tt_result.CpOrderId
	this.channelOrderId = this.tt_result.SDKOrderId
	payAmount := 0.0
	if payAmount, err = strconv.ParseFloat(this.tt_result.PayFee, 64); err != nil {
		return err
	} else {
		this.payAmount = int(payAmount * 100)
	}
	return
}

func (this *TT) ParseChannelRet() (err error) {
	if result := this.tt_result.PayResult; result != "1" {
		this.callbackRet = err_ttResultFailure
	}
	return
}

func (this *TT) parseBody() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_ttParseBody
			beego.Error(err)
		}
	}()

	var buffer bytes.Buffer
	if _, err = io.Copy(&buffer, this.ctx.Request.Body); err != nil {
		return
	}
	contentBody := string(buffer.Bytes())

	beego.Trace(contentBody)
	this.tt_contentBody, _ = url.QueryUnescape(contentBody)
	beego.Trace(this.tt_contentBody)
	if err = json.Unmarshal([]byte(this.tt_contentBody), &this.tt_result); err != nil {
		beego.Error(err)
		return err
	}

	this.tt_Sign = this.ctx.Request.Header.Get("Sign")
	beego.Trace("head OK:" + this.tt_Sign)
	return
}

func (this *TT) ParseParam() (err error) {
	if err = this.parseBody(); err != nil {
		return
	}
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

func (this *TT) CheckSign() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_checkSign
			beego.Error(err)
		}
	}()

	content := fmt.Sprintf("%s%s", this.tt_contentBody, this.payKey)
	beego.Trace("content:" + content)
	var result string
	if result, err = tool.TTSign(content); err != nil {
		beego.Error(err)
	}

	if result != this.tt_Sign {
		msg := fmt.Sprintf("Sign is invalid, sign:%s, urlSign:%s", result, this.tt_Sign)
		err = errors.New(msg)
		return
	}
	return
}

func (this *TT) GetResult() (ret string) {
	if this.callbackRet == err_noerror {
		ret = `{"head":{"result":"0","message":"成功"}}`
	} else {
		ret = `{"head":{"result":"1","message":"失败"}}`
	}
	return
}
