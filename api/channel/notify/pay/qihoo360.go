package channelPayNotify

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"net/url"
	"sort"
	"strconv"
	"u9/tool"
)

var qihoo360UrlKeys []string = []string{"app_key", "product_id", "amount", "app_uid",
	"user_id", "order_id", "gateway_flag", "sign_type", "sign_return", "sign"}

type Qihoo360 struct {
	Base
}

type qihoo360ResData struct {
	app_key      string `json:"app_key"`
	product_id   int    `json:"product_id"`
	amount       int    `json:"amount"`
	app_uid      string `json:"app_uid"`
	app_ext1     int    `json:"app_ext1"`
	app_ext2     string `json:"app_ext2"`
	user_id      int    `json:"user_id"`
	order_id     int    `json:"order_id"`
	gateway_flag string `json:"gateway_flag"`
	sign_type    string `json:"sign_type"`
	app_order_id string `json:"app_order_id"`
	sign_return  string `json:"sign_return"`
	sign         string `json:"sign"`
}

func NewQihoo360(channelId, productId int, urlParams *url.Values, ctx *context.Context) *Qihoo360 {
	ret := new(Qihoo360)
	ret.Init(channelId, productId, urlParams, ctx)
	return ret
}

func (this *Qihoo360) Init(channelId, productId int, urlParams *url.Values, ctx *context.Context) {
	this.Base.InitWithCtx(channelId, productId, urlParams, &qihoo360UrlKeys, ctx)
}

func (this *Qihoo360) parseUrlParam() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_parseUrlParam
			beego.Error(err)
			beego.Error(this.urlParams)
		}
	}()

	this.orderId = this.urlParams.Get("app_order_id")
	this.channelOrderId = this.urlParams.Get("order_id")
	this.channelUserId = this.urlParams.Get("user_id")

	payAmount := 0
	if payAmount, err = strconv.Atoi(this.urlParams.Get("amount")); err != nil {
		return err
	} else {
		this.payAmount = payAmount
	}
	return
}

func (this *Qihoo360) ParseParam() (err error) {
	if err = this.parseUrlParam(); err != nil {
		return
	}
	if err = this.parseChannelGameKey("QIHOO360_APP_KEY"); err != nil {
		return
	}
	if err = this.parseChannelPayKey("QIHOO360_SECRET"); err != nil {
		return
	}
	if err = this.Base.ParseParam(); err != nil {
		return
	}
	return
}

func (this *Qihoo360) CheckSign() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_checkSign
			beego.Error(err)
		}
	}()

	excludeItems := []string{"sign_return", "sign"}
	sorter := tool.NewUrlValuesSorter(this.urlParams, &excludeItems)
	sort.Sort(sorter)
	content := sorter.FormatBody("v", "#") + "#" + this.channelPayKey
	//beego.Trace(content)

	urlSign := this.urlParams.Get("sign")
	if sign := tool.Md5([]byte(content)); sign != urlSign {
		msg := fmt.Sprintf("Sign is invalid, content:%s, sign:%s, urlSign:%s", content, sign, urlSign)
		err = errors.New(msg)
		return
	}
	return
}

func (this *Qihoo360) GetResult() (ret string) {
	beego.Trace("callbackRet:" + strconv.Itoa(this.callbackRet))
	if this.callbackRet == err_noerror {
		ret = `ok`
	} else {
		ret = `fail`
	}
	return
}
