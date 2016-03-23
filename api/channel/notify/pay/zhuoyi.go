package channelPayNotify

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"net/url"
	"strconv"
	"u9/tool"
)

var zhuoyiUrlKeys []string = []string{"Recharge_Id", "App_Id", "Uin", "Urecharge_Id",
	"Recharge_Money", "Pay_Status", "Create_Time", "Sign"}

const (
	err_zhuoyiAppServerKey  = 12801
	err_zhuoyiResultFailure = 12802
)

type Zhuoyi struct {
	Base
	appServerKey string
}

func NewZhuoYi(channelId, productId int, urlParams *url.Values) *Zhuoyi {
	ret := new(Zhuoyi)
	ret.Init(channelId, productId, urlParams)
	return ret
}

func (this *Zhuoyi) Init(channelId, productId int, urlParams *url.Values) {
	this.Base.Init(channelId, productId, urlParams, &zhuoyiUrlKeys)
}

func (this *Zhuoyi) parseAppServerKey() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_zhuoyiAppServerKey
			beego.Trace(err)
		}
	}()
	this.appServerKey, err = this.getPackageParam("zy_app_secret")
	return
}

func (this *Zhuoyi) parseUrlParam() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_parseUrlParam
			beego.Trace(err)
		}
	}()

	this.orderId = this.urlParams.Get("Urecharge_Id")
	this.channelUserId = this.urlParams.Get("Uin")
	this.channelOrderId = this.urlParams.Get("Recharge_Id")

	payAmount := 0.0
	if payAmount, err = strconv.ParseFloat(this.urlParams.Get("Recharge_Money"), 64); err != nil {
		return err
	} else {
		this.payAmount = int(payAmount * 100)
	}
	return
}

func (this *Zhuoyi) ParseChannelRet() (err error) {
	if result := this.urlParams.Get("Pay_Status"); result != "1" {
		this.callbackRet = err_zhuoyiResultFailure
	}
	return
}

func (this *Zhuoyi) ParseParam() (err error) {
	if err = this.parseUrlParam(); err != nil {
		return
	}
	if err = this.parseAppServerKey(); err != nil {
		return
	}
	if err = this.Base.ParseParam(); err != nil {
		return
	}
	return
}

func (this *Zhuoyi) CheckSign() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_checkSign
			beego.Error(this.urlParams)
			beego.Error(err)
		}
	}()

	format := "App_Id=%s&Create_Time=%s&Extra=%s&Pay_Status=%s&Recharge_Gold_Count=%s&Recharge_Id=%s&Recharge_Money=%s&Uin=%s&Urecharge_Id=%s%s"
	content := fmt.Sprintf(format,
		this.urlParams.Get("App_Id"), this.urlParams.Get("Create_Time"),
		this.urlParams.Get("Extra"), this.urlParams.Get("Pay_Status"),
		this.urlParams.Get("Recharge_Gold_Count"), this.urlParams.Get("Recharge_Id"),
		this.urlParams.Get("Recharge_Money"), this.urlParams.Get("Uin"),
		this.urlParams.Get("Urecharge_Id"), this.appServerKey)

	urlSign := this.urlParams.Get("Sign")
	if sign := tool.Md5([]byte(content)); sign != urlSign {
		msg := fmt.Sprintf("Sign is invalid, content:%s, sign:%s, urlSign:%s", content, sign, urlSign)
		err = errors.New(msg)
		return
	}
	return
}

func (this *Zhuoyi) GetResult() (ret string) {
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
