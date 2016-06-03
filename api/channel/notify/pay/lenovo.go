package channelPayNotify

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"net/url"
	"strconv"
	"u9/tool"
)

type Lenovo struct {
	Base
}

type lenovoResData_TransData struct {
	ExOrderNo string `json:"exorderno"`
	TransId   string `json:"transid"`
	AppId     string `json:"appid"`
	WareSID   int    `json:"waresid"`
	FeeType   int    `json:"feetype"`
	Money     int    `json:"money"`
	Count     int    `json:"count"`
	Result    int    `json:"result"`
	TransType int    `json:"transtype"`
	TransTime string `json:"transtime"`
	CpPrivate string `json:"cpprivate"`
	PayType   int    `json:"paytype"`
}

var (
	lenovoRsaPrivateKey *rsa.PrivateKey
)

type lenovoResData struct {
	TransData lenovoResData_TransData `json:"transdata"`
	Sign      string                  `json:"sign"`
}

func NewLenovo(channelId, productId int, urlParams *url.Values, ctx *context.Context) *Lenovo {
	ret := new(Lenovo)
	ret.Init(channelId, productId, urlParams, ctx)
	return ret
}

func (this *Lenovo) Init(channelId, productId int, urlParams *url.Values, ctx *context.Context) {
	this.Base.InitWithCtx(channelId, productId, urlParams, &emptyUrlKeys, ctx)
	this.data = new(lenovoResData)
}

func (this *Lenovo) parseForm() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_htcParseBody
			beego.Error(this.ctx.Request.Form)
			beego.Error(err)
		}
	}()

	data := this.data.(*lenovoResData)
	data.Sign = this.ctx.Request.FormValue("sign")
	if err = json.Unmarshal([]byte(this.ctx.Request.FormValue("transdata")), &data.TransData); err != nil {
		return
	}

	this.orderId = data.TransData.ExOrderNo
	this.channelOrderId = data.TransData.TransId
	this.payAmount = data.TransData.Money

	beego.Trace(this.ctx.Request.Form)
	//beego.Trace(fmt.Sprintf("%+v", data))
	return
}

func (this *Lenovo) parseRsaPrivateKey() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_parseRsaPrivateKey
			beego.Error(err)
		}
	}()
	if lenovoRsaPrivateKey == nil {
		if err = this.parseChannelPayKey("lenovo.open.appkey"); err != nil {
			return
		}
		rsaPrivateKeyStr := this.channelPayKey
		if lenovoRsaPrivateKey, err = tool.ParsePkCS8PrivateKeyWithStr(rsaPrivateKeyStr); err != nil {
			return
		}
	}
	return nil
}

func (this *Lenovo) ParseParam() (err error) {
	if err = this.parseForm(); err != nil {
		return
	}
	if err = this.parseRsaPrivateKey(); err != nil {
		return
	}
	if err = this.Base.ParseParam(); err != nil {
		return
	}
	this.channelUserId = this.loginRequest.ChannelUserid
	return
}

func (this *Lenovo) CheckSign() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_checkSign
			beego.Error(err)
		}
	}()

	content := this.ctx.Request.FormValue("transdata")
	sign := ""
	if sign, err = tool.RsaPKCS1V15Sign(lenovoRsaPrivateKey, content); err != nil {
		return
	}

	data := this.data.(*lenovoResData)
	urlSign := data.Sign
	if sign != urlSign {
		if urlSign, err = url.QueryUnescape(urlSign); err != nil {
			return
		}
		if sign != urlSign {
			msg := fmt.Sprintf("Sign is invalid, content:%s, urlSign:%s", content, data.Sign)
			err = errors.New(msg)
		}
	}
	return
}

func (this *Lenovo) ParseChannelRet() (err error) {
	if err = this.Base.ParseChannelRet(); err != nil {
		return
	}

	data := this.data.(*lenovoResData)

	if err = this.parseChannelGameID("lenovo.open.appid"); err != nil {
		return
	}
	if this.channelGameId != data.TransData.AppId {
		this.callbackRet = err_parseChannelGameId
		beego.Error("lenovo.open.appid is invalid.")
		return
	}

	if data.TransData.Result != 0 {
		this.callbackRet = err_callbackFail
		beego.Error("transData.result is equal 0.")
		return
	}
	return
}

func (this *Lenovo) GetResult() (ret string) {
	beego.Trace("callbackRet:" + strconv.Itoa(this.callbackRet))
	if this.callbackRet == err_noerror {
		ret = `SUCCESS`
	} else {
		ret = `FAILURE`
	}
	return
}
