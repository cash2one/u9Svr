package channelPayNotify

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"net/url"
	"strconv"
	"u9/tool"
)

type Vivo struct {
	Base
}

type vivoResData struct {
	RespCode      string `json:"respCode"`
	RespMsg       string `json:"respMsg"`
	SignMethod    string `json:"signMethod"`
	Signature     string `json:"signature"`
	TradeType     string `json:"tradeType"`
	TradeStatus   string `json:"tradeStatus"`
	CpId          string `json:"cpId"`
	AppId         string `json:"appId"`
	CpOrderNumber string `json:"cpOrderNumber"`
	orderNumber   string `json:"orderNumber"`
	OrderAmount   string `json:"orderAmount"`
	ExtInfo       string `json:"extInfo"`
	PayTime       string `json:"payTime"`
	Uid           string `json:"uid"`
}

func NewVivo(channelId, productId int, urlParams *url.Values, ctx *context.Context) *Vivo {
	ret := new(Vivo)
	ret.Init(channelId, productId, urlParams, ctx)
	return ret
}

func (this *Vivo) Init(channelId, productId int, urlParams *url.Values, ctx *context.Context) {
	this.Base.InitWithCtx(channelId, productId, urlParams, &emptyUrlKeys, ctx)
	//this.data = new(vivoResData)
}

func (this *Vivo) parseUrlParam() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_parseUrlParam
			beego.Error(err)
			beego.Error(this.urlParams)
		}
	}()

	this.orderId = this.urlParams.Get("cpOrderNumber")
	this.channelUserId = this.urlParams.Get("uid")
	this.channelOrderId = this.urlParams.Get("orderNumber")

	payAmount := 0
	if payAmount, err = strconv.Atoi(this.urlParams.Get("orderAmount")); err != nil {
		return err
	} else {
		this.payAmount = payAmount
	}

	beego.Trace(this.urlParams)
	return
}

func (this *Vivo) ParseParam() (err error) {
	if err = this.parseUrlParam(); err != nil {
		return
	}
	if err = this.parseChannelGameKey("VIVO_CP_KEY"); err != nil {
		return
	}
	if err = this.Base.ParseParam(); err != nil {
		return
	}
	return
}

func (this *Vivo) CheckSign() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_checkSign
			beego.Error(err)
		}
	}()

	format := `appId=%s&cpId=%s&cpOrderNumber=%s&extInfo=%s&orderAmount=%s&orderNumber=%s&` +
		`payTime=%s&respCode=%s&respMsg=%s&tradeStatus=%s&tradeType=%s&uid=%s&%s`
	content := fmt.Sprintf(format,
		this.urlParams.Get("appId"), this.urlParams.Get("cpId"),
		this.urlParams.Get("cpOrderNumber"), this.urlParams.Get("extInfo"),
		this.urlParams.Get("orderAmount"), this.urlParams.Get("orderNumber"),
		this.urlParams.Get("payTime"), this.urlParams.Get("respCode"),
		this.urlParams.Get("respMsg"), this.urlParams.Get("tradeStatus"),
		this.urlParams.Get("tradeType"), this.urlParams.Get("uid"),
		tool.Md5([]byte(this.channelGameKey)))

	urlSign := this.urlParams.Get("signature")
	if sign := tool.Md5([]byte(content)); sign != urlSign {
		msg := fmt.Sprintf("Sign is invalid, content:%s, sign:%s, urlSign:%s", content, sign, urlSign)
		err = errors.New(msg)
		return
	}
	return
}

func (this *Vivo) GetResult() (ret string) {
	if this.callbackRet != err_noerror {
		//	this.ctx.Abort(403, msg)
		msg := "callbackRet:" + strconv.Itoa(this.callbackRet)
		this.ctx.ResponseWriter.WriteHeader(403)
		this.ctx.ResponseWriter.Write([]byte(msg))
		//panic(msg)
	}
	return
}
