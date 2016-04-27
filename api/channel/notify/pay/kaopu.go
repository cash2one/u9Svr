package channelPayNotify

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"net/url"
	"strconv"
	"u9/tool"
)

var kaopuUrlKeys []string = []string{"username", "kpordernum", "ywordernum", "status", "amount", "gamename", "sign"}

const (
	err_kaopuParsePayKey   = 10801
	err_kaopuResultFailure = 10802
)

//靠谱
type KaoPu struct {
	Base
	sign   string
	payKey string
}

func NewKaoPu(channelId, productId int, urlParams *url.Values) *KaoPu {
	ret := new(KaoPu)
	ret.Init(channelId, productId, urlParams)
	return ret
}

func (this *KaoPu) Init(channelId, productId int, urlParams *url.Values) {
	this.Base.Init(channelId, productId, urlParams, &kaopuUrlKeys)
}

func (this *KaoPu) parsePayKey() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_kaopuParsePayKey
			beego.Trace(err)
		}
	}()
	this.payKey, err = this.getPackageParam("KAOPU_SECRETKEY")
	return
}

func (this *KaoPu) parseUrlParam() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_parseUrlParam
			beego.Trace(err)
		}
	}()

	this.orderId = this.urlParams.Get("ywordernum")
	this.channelUserId = this.urlParams.Get("username")
	this.channelOrderId = this.urlParams.Get("kpordernum")
	var money string = this.urlParams.Get("amount")
	if this.payAmount, err = strconv.Atoi(money); err != nil {
		beego.Trace(err)
	}

	return
}

func (this *KaoPu) ParseChannelRet() (err error) {
	if result := this.urlParams.Get("status"); result != "1" {
		this.callbackRet = err_kaopuResultFailure
	}
	return
}

func (this *KaoPu) ParseParam() (err error) {
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

func (this *KaoPu) CheckSign() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_checkSign
			beego.Trace(err)
		}
	}()

	format := "%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s"
	content := fmt.Sprintf(format,
		this.channelUserId, this.channelOrderId, this.orderId, this.urlParams.Get("status"),
		this.urlParams.Get("paytype"), this.urlParams.Get("amount"), this.urlParams.Get("gameserver"),
		this.urlParams.Get("errdesc"), this.urlParams.Get("paytime"), this.urlParams.Get("gamename"), this.payKey)

	urlSign := this.urlParams.Get("sign")
	if this.sign = tool.Md5([]byte(content)); this.sign != urlSign {
		msg := fmt.Sprintf("Sign is invalid, content:%s, sign:%s, urlSign%s", content, this.sign, urlSign)
		err = errors.New(msg)
		return
	}

	return
}

// success_kaopu = `{"code":"1000","msg":"success","sign":"123"}`
// err_kaopuSign = `{"code":"1002","msg":"sign_erro","sign":"123"}`
// err_kaopuOreder = `{"code":"1003","msg":"order_erro","sign":"123"}`
// err_kaopuParam = `{"code":"1004","msg":"param_erro","sign":"123"}`
// err_kaopuSystem = `{"code":"1005","msg":"system_erro","sign":"123"}`
// err_kaopuUser = `{"code":"1006","msg":"user_erro","sign":"123"}`
// err_kaopuGameName = `{"code":"1007","msg":"gamename_erro","sign":"123"}`
// err_kaopuZone = `{"code":"1008","msg":"gameserver_erro","sign":"123"}`
// err_kaopuAmount = `{"code":"1009","msg":"amount_erro","sign":"123"}`
func (this *KaoPu) GetResult() (ret string) {
	type KaopuResultJson struct {
		Code string `json:"code"`
		Msg  string `json:"msg"`
		Sign string `json:"sign"`
	}
	kaopuResult := new(KaopuResultJson)
	switch this.callbackRet {
	case err_noerror:
		kaopuResult.Code = "1000"
		kaopuResult.Msg = "success"

	case err_checkSign:
		kaopuResult.Code = "1002"
		kaopuResult.Msg = "sign_err"
	case err_orderIsNotExist:
		kaopuResult.Code = "1003"
		kaopuResult.Msg = "order_err"
	case err_kaopuParsePayKey:
		kaopuResult.Code = "1004"
		kaopuResult.Msg = "param_err"
	case err_channelUserIsNotExist:
		kaopuResult.Code = "1006"
		kaopuResult.Msg = "user_err"
	case err_payAmountError:
		kaopuResult.Code = "1009"
		kaopuResult.Msg = "amount_err"
	default:
		kaopuResult.Code = "1005"
		kaopuResult.Msg = "system_err"
	}
	format := "%s|%s"
	content := fmt.Sprintf(format, kaopuResult.Code, this.payKey)
	kaopuResult.Sign = tool.Md5([]byte(content))
	data, _ := json.Marshal(kaopuResult)
	ret = string(data)
	beego.Trace(ret)
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
