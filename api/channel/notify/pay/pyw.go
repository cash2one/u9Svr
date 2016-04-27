package channelPayNotify

// import (
// 	"crypto/rsa"
// 	"encoding/json"
// 	// "errors"
// 	"bytes"
// 	"fmt"
// 	"github.com/astaxie/beego"
// 	"github.com/astaxie/beego/context"
// 	"io"
// 	"net/url"
// 	"strings"
// 	"u9/tool"
// )

// var pywUrlKeys []string = []string{"order", "sign", "sign_type"}

// const (
// 	err_PYWParsePayKey      = 12701
// 	err_PYWResultFailure    = 12702
// 	err_PYWInitRsaPublicKey = 12703
// 	err_PYWParseBody        = 12704
// )

// //PYW
// type PYW struct {
// 	Base
// 	payKey     string
// 	response   Response
// 	PYW_result PYW_Result
// 	ctx        *context.Context
// }

// type PYW_Result struct {
// 	Result_code   int    `json:"result_code"`
// 	Gmt_create    string `json:"gmt_create"`
// 	Real_amount   int    `json:"real_amount"`
// 	Result_msg    string `json:"result_msg"`
// 	Game_code     string `json:"game_code"`
// 	Game_order_id string `json:"game_order_id"`
// 	Jolo_order_id string `json:"jolo_order_id"`
// 	Gmt_payment   string `json:"gmt_payment"`
// }

// var (
// 	PYWRsaPublicKey *rsa.PublicKey
// )

// func NewPYW(channelId, productId int, urlParams *url.Values, ctx *context.Context) *PYW {
// 	ret := new(PYW)
// 	ret.Init(channelId, productId, urlParams, ctx)
// 	return ret
// }

// func (this *PYW) Init(channelId, productId int, urlParams *url.Values, ctx *context.Context) {
// 	this.Base.Init(channelId, productId, urlParams, &pywUrlKeys)
// 	this.ctx = ctx
// }

// func (this *PYW) parsePayKey() (err error) {
// 	defer func() {
// 		if err != nil {
// 			this.callbackRet = err_PYWParsePayKey
// 			beego.Trace(err)
// 		}
// 	}()
// 	this.payKey, err = this.getPackageParam("PYW_SDK_PUBLICKEY")
// 	return
// }

// func (this *PYW) CheckUrlParam() (err error) {
// 	return
// }

// func (this *PYW) parseUrlParam() (err error) {
// 	defer func() {
// 		if err != nil {
// 			this.callbackRet = err_parseUrlParam
// 			beego.Trace(err)
// 		}
// 	}()
// 	// this.order = url.QueryEscape(this.response.Order)

// 	this.response.Order = strings.Replace(this.response.Order, "\"{", "{", 1)
// 	this.response.Order = strings.Replace(this.response.Order, "}\"", "}", 1)
// 	beego.Trace(this.response.Order)
// 	json.Unmarshal([]byte(this.response.Order), &this.PYW_result)
// 	beego.Trace(this.PYW_result)
// 	this.orderId = this.PYW_result.Game_order_id
// 	this.channelOrderId = this.PYW_result.Jolo_order_id
// 	this.payAmount = this.PYW_result.Real_amount

// 	return
// }

// func (this *PYW) ParseChannelRet() (err error) {
// 	if result := this.PYW_result.Result_code; result != 1 {
// 		this.callbackRet = err_PYWResultFailure
// 	}
// 	return
// }
// func (this *PYW) parseBody() (err error) {
// 	defer func() {
// 		if err != nil {
// 			this.callbackRet = err_PYWParseBody
// 			beego.Trace(err)
// 		}
// 	}()

// 	var buffer bytes.Buffer
// 	if _, err = io.Copy(&buffer, this.ctx.Request.Body); err != nil {
// 		return
// 	}
// 	content := string(buffer.Bytes())
// 	var newValues url.Values
// 	if newValues, err = url.ParseQuery(content); err != nil {
// 		return
// 	}
// 	this.response.Order = newValues.Get("order")
// 	this.response.Sign = newValues.Get("sign")
// 	this.response.Sign_type = newValues.Get("sign_type")
// 	this.response.Sign = strings.Replace(this.response.Sign, "\"", "", 2)
// 	beego.Trace(this.response.Order)
// 	beego.Trace(this.response.Sign)
// 	beego.Trace(this.response.Sign_type)
// 	return
// }

// func (this *PYW) ParseParam() (err error) {
// 	if err = this.parseBody(); err != nil {
// 		return
// 	}
// 	if err = this.parseUrlParam(); err != nil {
// 		return
// 	}
// 	if err = this.parsePayKey(); err != nil {
// 		return
// 	}
// 	if err = this.Base.ParseParam(); err != nil {
// 		return
// 	}
// 	this.channelUserId = this.loginRequest.ChannelUserid
// 	this.initRsaPublicKey()
// 	return
// }

// func (this *PYW) initRsaPublicKey() (err error) {
// 	defer func() {
// 		if err != nil {
// 			this.callbackRet = err_PYWInitRsaPublicKey
// 			beego.Trace(err)
// 		}
// 	}()

// 	if PYWRsaPublicKey == nil {
// 		PYWRsaPublicKey, err = tool.ParsePKIXPublicKeyWithStr(this.payKey)
// 		if err != nil {
// 			beego.Error(err)
// 			return err
// 		}
// 	}
// 	// beego.Trace(PYWRsaPublicKey)
// 	return nil
// }

// func (this *PYW) CheckSign() (err error) {
// 	defer func() {
// 		if err != nil {
// 			this.callbackRet = err_checkSign
// 			beego.Trace(err)
// 		}
// 	}()

// 	// if sign := tool.RsaVerifyPKCS1v15(PYWRsaPublicKey, this.order); sign != this.urlParams.Get("signature") {
// 	// 	msg := fmt.Sprintf("Sign is invalid, context:%s, sign:%s", context, sign)
// 	// 	err = errors.New(msg)
// 	// 	return
// 	// }
// 	if err = tool.RsaVerifyPKCS1v15(PYWRsaPublicKey, this.response.Order, this.response.Sign); err != nil {
// 		msg := fmt.Sprintf("RsaVerifyPK CS1v15 exception: context:%s, sign:%s", this.response.Order, this.response.Sign)
// 		beego.Trace(msg)
// 		return err
// 	}
// 	return
// }

// func (this *PYW) GetResult() (ret string) {
// 	if this.callbackRet == err_noerror {
// 		ret = "success"
// 	} else {
// 		ret = "failure"
// 	}
// 	return
// }

// /*
//   signature rule: md5("order=xxxx&money=xxxx&mid=xxxx&time=xxxx&result=x&ext=xxx&key=xxxx")
//   test url:
//   http://192.168.0.185/api/channelPayNotify/1000/101/?
//   order=test20160116172500359&
//   money=100.00&
//   mid=test10086001&
//   time=20160116172500&
//   result=1&
//   ext=game20160116175128772&
//   signature=8f00a109716e819bfe0afb695c1addf1
// */
