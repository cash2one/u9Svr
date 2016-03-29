package channelPayNotify

import (
	"encoding/json"
	// "errors"
	"github.com/astaxie/beego"
	"net/url"
	// "strconv"
	"strings"
	"u9/tool"
)

var yyhUrlKeys []string = []string{"transdata", "sign"}

const (
	err_yyhParsePayKey   = 10101
	err_yyhResultFailure = 10102
)

//应用汇
type YYH struct {
	Base
	payKey    string
	transData TransData
}
type TransData struct {
	Exorderno string `json:"exorderno"`
	Transid   string `json:"transid"`
	Appid     string `json:"appid"`
	Waresid   int    `json:"waresid"`
	Feetype   int    `json:"feetype"`
	Money     int    `json:"money"`
	Count     int    `json:"count"`
	Result    int    `json:"result`
	Transtype int    `json:"transtype"`
	Transtime string `json:"transtime"`
	Cpprivate string `json:"cpprivate"`
	PayType   string `json:"paytype"`
}

func NewYYH(channelId, productId int, urlParams *url.Values) *YYH {
	ret := new(YYH)
	ret.Init(channelId, productId, urlParams)
	return ret
}

func (this *YYH) Init(channelId, productId int, urlParams *url.Values) {
	this.Base.Init(channelId, productId, urlParams, &yyhUrlKeys)
}

func (this *YYH) parsePayKey() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_yyhParsePayKey
			beego.Trace(err)
		}
	}()
	this.payKey, err = this.getPackageParam("YYH_APPKEY")
	return
}

func (this *YYH) parseUrlParam() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_parseUrlParam
			beego.Trace(err)
		}
	}()
	beego.Trace(this.urlParams)
	json.Unmarshal([]byte(this.urlParams.Get("transdata")), &this.transData)
	this.orderId = this.transData.Exorderno
	this.channelOrderId = this.transData.Transid
	this.payAmount = this.transData.Money

	return
}

func (this *YYH) ParseChannelRet() (err error) {
	if result := this.transData.Result; result != 0 {
		this.callbackRet = err_yyhResultFailure
		beego.Trace(result)
		return
	}
	return
}

func (this *YYH) ParseParam() (err error) {
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
	return
}

func (this *YYH) CheckSign() (err error) {
	var result string
	md5Sign := tool.Md5([]byte(this.urlParams.Get("transdata")))
	if result, err = tool.YYHSign(md5Sign, this.urlParams.Get("sign"), this.payKey); err != nil {
		beego.Trace(err)
		return
	}
	result = strings.TrimSpace(result)
	if result != "0" {
		this.callbackRet = err_checkSign
		beego.Trace("yyh check:", result, "transdata:", md5Sign, "sign:", this.urlParams.Get("sign"), "paykey:", this.payKey)
	} else {
		beego.Trace("yyh check:", result)
	}

	return
}

func (this *YYH) GetResult() (ret string) {
	if this.callbackRet == err_noerror {
		ret = "true"
	} else {
		ret = "false"
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
