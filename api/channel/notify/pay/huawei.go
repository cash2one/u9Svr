package channelPayNotify

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"net/url"
	"sort"
	"strconv"
	"u9/tool"
)

type Huawei struct {
	Base
}

var (
	huaweiRsaPublicKey *rsa.PublicKey
)

type huaweiResData struct {
	Result      string `json:"result"`
	UserName    string `json:"userName"`
	ProductName string `json:"productName"`
	PayType     string `json:"payType"`
	Amount      string `json:"amount"`
	OrderId     string `json:"orderId"`
	NotifyTime  string `json:"notifyTime"`
	RequestId   string `json:"requestId"`
	BankId      string `json:"bankId"`
	OrderTime   string `json:"orderTime"`
	TradeTime   string `json:"tradeTime"`
	AccessMode  string `json:"accessMode"`
	Spending    string `json:"spending"`
	ExtReserved string `json:"extReserved"`
	Sign        string `json:"sign"`
}

func NewHuawei(channelId, productId int, urlParams *url.Values, ctx *context.Context) *Huawei {
	ret := new(Huawei)
	ret.Init(channelId, productId, urlParams, ctx)
	return ret
}

func (this *Huawei) Init(channelId, productId int, urlParams *url.Values, ctx *context.Context) {
	this.Base.InitWithCtx(channelId, productId, urlParams, &emptyUrlKeys, ctx)
}

func (this *Huawei) parseUrlParam() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_parseUrlParam
			beego.Error(err)
			beego.Error(fmt.Sprintf("%+v", this.urlParams))
		}
	}()

	this.orderId = this.urlParams.Get("requestId")
	this.channelUserId = this.urlParams.Get("userName")
	this.channelOrderId = this.urlParams.Get("orderId")

	payAmount := 0.0
	if payAmount, err = strconv.ParseFloat(this.urlParams.Get("amount"), 64); err != nil {
		return err
	} else {
		this.payAmount = int(payAmount * 100)
	}
	return
}

func (this *Huawei) parseRsaPublicKey() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_parseRsaPublicKey
			beego.Error(err)
		}
	}()
	if huaweiRsaPublicKey == nil {
		rsaPublicKeyStr := ""
		if rsaPublicKeyStr, err = this.getPackageParam("HUAWEI_PAY_PUBLIC_KEY"); err != nil {
			return
		}
		//devPublicKey test
		//rsaPublicKeyStr = `MFwwDQYJKoZIhvcNAQEBBQADSwAwSAJBAIW1g+KAqqOeC1ypte8L3qTDk2nz6jUbM6o6Jg9obvivPnCAm/wZvV3jWbYWfOuO/wrFJygn/jZqf8cR1T1CQa8CAwEAAQ==`
		if huaweiRsaPublicKey, err = tool.ParsePKIXPublicKeyWithStr(rsaPublicKeyStr); err != nil {
			return
		}
	}
	return nil
}

func (this *Huawei) ParseParam() (err error) {
	if err = this.parseUrlParam(); err != nil {
		return
	}
	if err = this.parseRsaPublicKey(); err != nil {
		return
	}
	if err = this.Base.ParseParam(); err != nil {
		return
	}
	err = nil
	return
}

func (this *Huawei) ParseChannelRet() (err error) {
	if result := this.urlParams.Get("result"); result != "0" {
		this.callbackRet = err_callbackFail
	}
	return
}

func (this *Huawei) CheckSign() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_checkSign
			beego.Error(err)
		}
	}()

	excludeItems := []string{"sign"}
	sorter := tool.NewUrlValuesSorter(this.urlParams, &excludeItems)
	sort.Sort(sorter)
	content := sorter.Body()
	beego.Trace(content)
	urlSign := this.urlParams.Get("sign")
	if err = tool.RsaVerifyPKCS1v15(huaweiRsaPublicKey, content, urlSign); err != nil {
		msg := fmt.Sprintf("Sign is invalid, content:%s, urlSign:%s", content, urlSign)
		err = errors.New(msg)
	}

	return
}

func (this *Huawei) GetResult() (ret string) {
	beego.Trace("callbackRet:" + strconv.Itoa(this.callbackRet))
	if this.callbackRet == err_noerror {
		ret = `SUCCESS`
	} else {
		ret = `FAILURE`
	}
	return
}
