package channelPayNotify

import (
	"crypto/rsa"
	"encoding/json"
	// "errors"
	"bytes"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"io"
	"net/url"
	"strings"
	"u9/tool"
)

var htcUrlKeys []string = []string{"order", "sign", "sign_type"}

const (
	err_htcParsePayKey      = 12701
	err_htcResultFailure    = 12702
	err_htcInitRsaPublicKey = 12703
	err_htcParseBody        = 12704
)

//HTC
type HTC struct {
	Base
	payKey     string
	response   Response
	htc_result HTC_Result
	ctx        *context.Context
}
type Response struct {
	Order     string
	Sign      string
	Sign_type string
}
type HTC_Result struct {
	Result_code   int    `json:"result_code"`
	Gmt_create    string `json:"gmt_create"`
	Real_amount   int    `json:"real_amount"`
	Result_msg    string `json:"result_msg"`
	Game_code     string `json:"game_code"`
	Game_order_id string `json:"game_order_id"`
	Jolo_order_id string `json:"jolo_order_id"`
	Gmt_payment   string `json:"gmt_payment"`
}

var (
	htcRsaPublicKey *rsa.PublicKey
)

func NewHTC(channelId, productId int, urlParams *url.Values, ctx *context.Context) *HTC {
	ret := new(HTC)
	ret.Init(channelId, productId, urlParams, ctx)
	return ret
}

func (this *HTC) Init(channelId, productId int, urlParams *url.Values, ctx *context.Context) {
	this.Base.Init(channelId, productId, urlParams, &htcUrlKeys)
	this.ctx = ctx
}

func (this *HTC) parsePayKey() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_htcParsePayKey
			beego.Trace(err)
		}
	}()
	this.payKey, err = this.getPackageParam("HTC_SDK_PUBLICKEY")
	return
}

func (this *HTC) CheckUrlParam() (err error) {
	return
}

func (this *HTC) parseUrlParam() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_parseUrlParam
			beego.Error(err)
		}
	}()
	// this.order = url.QueryEscape(this.response.Order)

	this.response.Order = strings.Replace(this.response.Order, "\"{", "{", 1)
	this.response.Order = strings.Replace(this.response.Order, "}\"", "}", 1)
	beego.Trace(this.response.Order)
	json.Unmarshal([]byte(this.response.Order), &this.htc_result)
	beego.Trace(this.htc_result)
	this.orderId = this.htc_result.Game_order_id
	this.channelOrderId = this.htc_result.Jolo_order_id
	this.payAmount = this.htc_result.Real_amount

	return
}

func (this *HTC) ParseChannelRet() (err error) {
	if result := this.htc_result.Result_code; result != 1 {
		this.callbackRet = err_htcResultFailure
	}
	return
}
func (this *HTC) parseBody() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_htcParseBody
			beego.Error(err)
		}
	}()

	var buffer bytes.Buffer
	if _, err = io.Copy(&buffer, this.ctx.Request.Body); err != nil {
		return
	}
	content := string(buffer.Bytes())
	var newValues url.Values
	if newValues, err = url.ParseQuery(content); err != nil {
		return
	}
	this.response.Order = newValues.Get("order")
	this.response.Sign = newValues.Get("sign")
	this.response.Sign_type = newValues.Get("sign_type")
	this.response.Sign = strings.Replace(this.response.Sign, "\"", "", 2)
	beego.Trace(this.response.Order)
	beego.Trace(this.response.Sign)
	beego.Trace(this.response.Sign_type)
	return
}

func (this *HTC) ParseParam() (err error) {
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
	this.channelUserId = this.loginRequest.ChannelUserid
	this.initRsaPublicKey()
	return
}

func (this *HTC) initRsaPublicKey() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_htcInitRsaPublicKey
			beego.Trace(err)
		}
	}()

	if htcRsaPublicKey == nil {
		htcRsaPublicKey, err = tool.ParsePKIXPublicKeyWithStr(this.payKey)
		if err != nil {
			beego.Error(err)
			return err
		}
	}
	// beego.Trace(htcRsaPublicKey)
	return nil
}

func (this *HTC) CheckSign() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_checkSign
			beego.Trace(err)
		}
	}()

	// if sign := tool.RsaVerifyPKCS1v15(htcRsaPublicKey, this.order); sign != this.urlParams.Get("signature") {
	// 	msg := fmt.Sprintf("Sign is invalid, context:%s, sign:%s", context, sign)
	// 	err = errors.New(msg)
	// 	return
	// }
	if err = tool.RsaVerifyPKCS1v15(htcRsaPublicKey, this.response.Order, this.response.Sign); err != nil {
		msg := fmt.Sprintf("RsaVerifyPK CS1v15 exception: context:%s, sign:%s", this.response.Order, this.response.Sign)
		beego.Trace(msg)
		return err
	}
	return
}

func (this *HTC) GetResult() (ret string) {
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
