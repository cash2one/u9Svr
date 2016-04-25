package channelPayNotify

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"net/url"
	"strconv"
	"u9/tool"
)

var snailUrlKeys []string = []string{"AppId", "ConsumeStreamId", "CooOrderSerial", "Uin", "OriginalMoney", "OrderMoney", "PayStatus"}

const (
	err_snailParsePayKey   = 13101
	err_snailResultFailure = 13102
)

//免商店
type Snail struct {
	Base
	appid string
	payKey string
}

func NewSnail(channelId, productId int, urlParams *url.Values) *Snail {
	ret := new(Snail)
	ret.Init(channelId, productId, urlParams)
	return ret
}

func (this *Snail) Init(channelId, productId int, urlParams *url.Values) {
	this.Base.Init(channelId, productId, urlParams, &snailUrlKeys)
}

func (this *Snail) parsePayKey() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_snailParsePayKey
			beego.Trace(err)
		}
	}()
	this.appid, err = this.getPackageParam("SNAIL_APPID")
	this.payKey, err = this.getPackageParam("SNAIL_APPKEY")
	return
}

func (this *Snail) parseUrlParam() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_parseUrlParam
			beego.Trace(err)
		}
	}()

	this.orderId = this.urlParams.Get("CooOrderSerial")
	this.channelUserId = this.urlParams.Get("Uin")
	this.channelOrderId = this.urlParams.Get("ConsumeStreamId")

	payAmount := 0.0
	if payAmount, err = strconv.ParseFloat(this.urlParams.Get("OriginalMoney"), 64); err != nil {
		return err
	} else {
		this.payAmount = int(payAmount * 100)
	}
	return
}

func (this *Snail) ParseChannelRet() (err error) {
	if result := this.urlParams.Get("PayStatus"); result != "1" {
		this.callbackRet = err_snailResultFailure
	}
	return
}

func (this *Snail) ParseParam() (err error) {
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

func (this *Snail) CheckSign() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_checkSign
			beego.Trace(err)
		}
	}()

	format := "%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s"
	content := fmt.Sprintf(format,
		this.appid,this.urlParams.Get("Act"), this.urlParams.Get("ProductName"),this.channelOrderId,
		this.orderId,this.channelUserId, this.urlParams.Get("GoodsId"), this.urlParams.Get("GoodsInfo"),
		this.urlParams.Get("GoodsCount"),this.urlParams.Get("OriginalMoney"),this.urlParams.Get("OrderMoney"),
		this.urlParams.Get("Note"),this.urlParams.Get("PayStatus"),this.urlParams.Get("CreateTime"),this.payKey)

	urlSign := this.urlParams.Get("Sign")
	if sign := tool.Md5([]byte(content)); sign != urlSign {
		msg := fmt.Sprintf("Sign is invalid, content:%s, sign:%s, urlSign:%s", content, sign, urlSign)
		err = errors.New(msg)
		return
	}

	return
}

func (this *Snail) GetResult() (ret string) {
	switch this.callbackRet {
	case err_noerror:
		ret = `{"ErrorCode":"1","ErrprDesc":"接受成功"}`
	case err_checkSign:
		ret = `{"ErrorCode":"5","ErrorDesc":"Sign 无效"}`
	case err_snailParsePayKey:
		ret = `{"ErrorCode":"2","ErrorDesc":"AppId 无效"}`
	default:
		ret = `{"ErrorCode":"0","ErrorDesc":"接受失败"}`
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
