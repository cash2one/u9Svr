package channelPayNotify

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"net/url"
	"strconv"
	// "u9/models"
	"u9/tool"
)

var AnZhiUrlKeys []string = []string{"data"}

const (
	err_anzhiParseAppSecret = 10301
)

//安智
type AnZhi struct {
	Base
	appSecret  string
	data       string
	paySuccess bool
	anZhiJson  anZhiData
}

type anZhiData struct {
	PayAmount    string `json:"payAmount"`    //支付金额
	Uid          string `json:"uid"`          //用户id
	NotifyTime   int    `json:"notifyTime"`   //请求时间
	CpInfo       string `json:"cpInfo"`       //回调信息
	Memo         string `json:"memo"`         //备注
	OrderAmount  string `json:"orderAmount"`  //订单金额
	OrderAccount string `json:"orderAccount"` //订单数量
	Code         string `json:"code"`         //订单状态
	OrderTime    string `json:"orderTime"`    //订单时间
	Msg          string `json:"msg"`          //消息
	OrderId      string `json:"orderId"`      //订单号
}

func NewAnZhi(channelId, productId int, urlParams *url.Values) *AnZhi {
	ret := new(AnZhi)
	ret.Init(channelId, productId, urlParams)
	return ret
}

func (this *AnZhi) Init(channelId, productId int, urlParams *url.Values) {
	this.Base.Init(channelId, productId, urlParams, &AnZhiUrlKeys)
}

func (this *AnZhi) parseAppSecret() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_anzhiParseAppSecret
			beego.Trace(err)
		}
	}()
	if this.appSecret, err = this.getPackageParam("ANZHI_APPSECRET"); err != nil {
		return
	}
	beego.Trace(this.appSecret)
	return
}

func (this *AnZhi) parseUrlParam() (err error) {
	this.data = this.urlParams.Get("data")
	beego.Trace(this.data)
	ret, erro := tool.JavaDesDecyrpt(this.appSecret, this.data)
	if erro != nil {
		beego.Trace(erro)
	}

	beego.Trace(ret)
	err = json.Unmarshal([]byte(ret), &this.anZhiJson)
	if err != nil {
		beego.Trace(err)
	}
	this.orderId = this.anZhiJson.CpInfo
	this.channelUserId = this.anZhiJson.Uid
	this.channelOrderId = this.anZhiJson.OrderId
	beego.Trace(this.anZhiJson.OrderAmount)
	payAmount := 0.0
	if payAmount, err = strconv.ParseFloat(this.anZhiJson.OrderAmount, 64); err != nil {
		beego.Trace(err)
		return err
	} else {
		this.payAmount = int(payAmount)
	}
	return
}

func (this *AnZhi) ParseParam() (err error) {
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

func (this *AnZhi) CheckSign() (err error) {
	return
}

func (this *AnZhi) GetResult() (ret string) {
	beego.Trace(this.callbackRet)
	if this.callbackRet == err_noerror {
		ret = "success"
	} else {
		ret = "failure"
	}
	return
}
