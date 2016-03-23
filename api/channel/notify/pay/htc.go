package channelPayNotify

import (
	"crypto/rsa"
	"encoding/json"
	// "errors"
	"fmt"
	"github.com/astaxie/beego"
	"net/url"
	"strconv"
	"u9/tool"
)

var htcUrlKeys []string = []string{"order", "sign"}

const (
	err_htcParsePayKey      = 12701
	err_htcResultFailure    = 12702
	err_htcInitRsaPublicKey = 12703
)

//HTC
type HTC struct {
	Base
	order      string
	payKey     string
	htc_result HTC_Result
}
type HTC_Result struct {
	Result_code   string `json:"result_code"`
	Gmt_create    string `json:"gmt_create"`
	Real_amount   string `json:"real_amount"`
	Result_msg    string `json:"result_msg"`
	Game_code     string `json:"game_code"`
	Game_order_id string `json:"game_order_id"`
	Jolo_order_id string `json:"jolo_order_id"`
	Gmt_payment   string `json:"gmt_payment"`
}

var (
	htcRsaPublicKey *rsa.PublicKey
)

func NewHTC(channelId, productId int, urlParams *url.Values) *HTC {
	ret := new(HTC)
	ret.Init(channelId, productId, urlParams)
	return ret
}

func (this *HTC) Init(channelId, productId int, urlParams *url.Values) {
	this.Base.Init(channelId, productId, urlParams, &htcUrlKeys)
}

func (this *HTC) parsePayKey() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_htcParsePayKey
			beego.Trace(err)
		}
	}()
	this.payKey, err = this.getPackageParam("HTC_PRIVATE_KEY")
	return
}

func (this *HTC) parseUrlParam() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_parseUrlParam
			beego.Trace(err)
		}
	}()
	this.order = url.QueryEscape(this.urlParams.Get("order"))
	json.Unmarshal([]byte(this.order), &this.htc_result)
	this.orderId = this.htc_result.Jolo_order_id
	this.channelOrderId = this.htc_result.Game_order_id

	payAmount := 0.0
	if payAmount, err = strconv.ParseFloat(this.htc_result.Gmt_payment, 64); err != nil {
		return err
	} else {
		this.payAmount = int(payAmount)
	}
	return
}

func (this *HTC) ParseChannelRet() (err error) {
	if result := this.htc_result.Result_code; result != "1" {
		this.callbackRet = err_htcResultFailure
	}
	return
}

func (this *HTC) ParseParam() (err error) {
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
	sign := this.urlParams.Get("sign")
	if err = tool.RsaVerifyPKCS1v15(htcRsaPublicKey, this.order, sign); err != nil {
		msg := fmt.Sprintf("RsaVerifyPK CS1v15 exception: context:%s, sign:%s", this.order, sign)
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
