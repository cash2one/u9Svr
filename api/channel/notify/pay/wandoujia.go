package channelPayNotify

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"net/url"
	"strconv"
	"u9/tool"
)

var wandoujiaUrlKeys []string = []string{"content", "signType", "sign"}

const (
	err_wandoujiaParseSecrectKey  = 11701
	err_wandoujiaInitRsaPublicKey = 11702
	wandoujiaRsaPublicKeyStr      = `MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCd95FnJFhPinpNiE/h4VA6bU1rzRa5+a25BxsnFX8TzquWxqDCoe4xG6QKXMXuKvV57tTRpzRo2jeto40eHKClzEgjx9lTYVb2RFHHFWio/YGTfnqIPTVpi7d7uHY+0FZ0lYL5LlW4E2+CQMxFOPRwfqGzMjs1SDlH7lVrLEVy6QIDAQAB`
)

var (
	wandoujiaRsaPublicKey *rsa.PublicKey
)

type WandoujiaLrRet struct {
	TimeStamp  uint64 `json:"timeStamp"`
	OrderId    uint64 `json:"orderId"`
	Money      int    `json:"money"`
	ChargeType string `json:"chargeType"`
	AppKeyId   uint64 `json:"appKeyId"`
	BuyerId    uint64 `json:"buyerId"`
	OutTradeNo string `json:"out_trade_no"`
	CardNo     uint64 `json:"cardNo"`
}

type Wandoujia struct {
	Base
	secrectKey string
	channelRet WandoujiaLrRet
}

func NewWandoujia(channelId, productId int, urlParams *url.Values) *Wandoujia {
	ret := new(Wandoujia)
	ret.Init(channelId, productId, urlParams)
	return ret
}

func (this *Wandoujia) Init(channelId, productId int, urlParams *url.Values) {
	this.Base.Init(channelId, productId, urlParams, &wandoujiaUrlKeys)
}

func (this *Wandoujia) initRsaPublicKey() (err error) {
	if wandoujiaRsaPublicKey == nil {
		wandoujiaRsaPublicKey, err = tool.ParsePKIXPublicKeyWithStr(wandoujiaRsaPublicKeyStr)
		if err != nil {
			this.callbackRet = err_wandoujiaInitRsaPublicKey
			beego.Error(err)
			return err
		}
	}
	return nil
}

func (this *Wandoujia) parseSecrectKey() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_wandoujiaParseSecrectKey
			beego.Trace(err)
		}
	}()
	this.secrectKey, err = this.getPackageParam("WANDOUJIA_SECRETKEY")
	return
}

func (this *Wandoujia) parseUrlParam() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_parseUrlParam
			beego.Trace(err)
		}
	}()

	content := this.urlParams.Get("content")
	if err = json.Unmarshal([]byte(content), &this.channelRet); err != nil {
		beego.Trace(content)
		return
	}

	this.orderId = this.channelRet.OutTradeNo
	this.channelOrderId = strconv.FormatUint(this.channelRet.OrderId, 10)
	this.channelUserId = strconv.FormatUint(this.channelRet.BuyerId, 10)

	this.payAmount = this.channelRet.Money

	return
}

func (this *Wandoujia) ParseParam() (err error) {
	if err = this.initRsaPublicKey(); err != nil {
		return
	}
	if err = this.parseSecrectKey(); err != nil {
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

func (this *Wandoujia) CheckSign() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_checkSign
			beego.Trace(err)
		}
	}()

	sign := this.urlParams.Get("sign")
	content := this.urlParams.Get("content")

	if err = tool.RsaVerifyPKCS1v15(wandoujiaRsaPublicKey, content, sign); err != nil {
		msg := fmt.Sprintf("RsaVerifyPK CS1v15 exception: context:%s, sign:%s", content, sign)
		beego.Trace(msg)
		return err
	}
	return
}

func (this *Wandoujia) GetResult() (ret string) {
	if this.callbackRet == err_noerror {
		ret = "success"
	} else {
		ret = "fail"
	}
	return
}
