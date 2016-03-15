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

var m4399UrlKeys []string = []string{"orderid", "p_type", "uid", "money", "gamemoney",
	"mark", "time", "sign"}

const (
	err_m4399ParseSecrectKey = 10901
)

//4399
type M4399 struct {
	Base
	secrectKey string
}

func NewM4399(channelId, productId int, urlParams *url.Values) *M4399 {
	ret := new(M4399)
	ret.Init(channelId, productId, urlParams)
	return ret
}

func (this *M4399) Init(channelId, productId int, urlParams *url.Values) {
	this.Base.Init(channelId, productId, urlParams, &m4399UrlKeys)
}

func (this *M4399) parseSecrectKey() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_m4399ParseSecrectKey
			beego.Trace(err)
		}
	}()
	this.secrectKey, err = this.getPackageParam("M4399_SECRECT")
	return
}

func (this *M4399) parseUrlParam() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_parseUrlParam
			beego.Trace(err)
		}
	}()

	this.orderId = this.urlParams.Get("mark")
	this.channelUserId = this.urlParams.Get("uid")
	this.channelOrderId = this.urlParams.Get("orderid")

	payAmount := 0.0
	if payAmount, err = strconv.ParseFloat(this.urlParams.Get("money"), 64); err != nil {
		return err
	} else {
		this.payAmount = int(payAmount * 100)
	}
	return
}

func (this *M4399) ParseParam() (err error) {
	if err = this.parseUrlParam(); err != nil {
		return
	}
	if err = this.parseSecrectKey(); err != nil {
		return
	}
	if err = this.Base.ParseParam(); err != nil {
		return
	}
	return
}

func (this *M4399) CheckSign() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_checkSign
			beego.Trace(err)
		}
	}()

	context := fmt.Sprintf("%s%s%s%s", this.channelOrderId,
		this.channelId, this.urlParams.Get("money"), this.urlParams.Get("gamemoney"))
	if this.urlParams.Get("serverid") != "" {
		context = context + this.urlParams.Get("serverid")
	}
	context = context + this.secrectKey + this.orderId
	if this.urlParams.Get("roleid") != "" {
		context = context + this.urlParams.Get("roleid")
	}
	context = context + this.urlParams.Get("time")

	if sign := tool.Md5([]byte(context)); sign != this.urlParams.Get("sign") {

		msg := fmt.Sprintf("Sign is invalid, context:%s, sign:%s", context, sign)
		err = errors.New(msg)
		beego.Trace(err)
		return
	}
	return
}

func (this *M4399) GetResult() (ret string) {
	type m4399Ret struct {
		Status    int    `json:"status"` //1：异常；2：成功；3：失败（将钱返还给用户）
		Code      string `json:"code"`
		Money     string `json:"money"`
		Gamemoney string `json:"gamemoney"`
		Msg       string `json:"msg"`
	}
	jsonRet := m4399Ret{
		Status:    1,
		Code:      "",
		Money:     "0",
		Gamemoney: "0",
		Msg:       "",
	}

	switch this.callbackRet {
	case err_noerror:
		jsonRet.Status = 2
		jsonRet.Msg = "success"
	case err_checkSign:
		jsonRet.Status = 1
		jsonRet.Msg = "sign_error"
	case err_parseLoginRequest:
		fallthrough
	case err_channelUserIsNotExist:
		jsonRet.Status = 1
		jsonRet.Msg = "user_not_exist"
	case err_handleOrder:
		jsonRet.Status = 1
		jsonRet.Msg = "orderid_exist"
	case err_payAmountError:
		jsonRet.Status = 1
		jsonRet.Msg = "money_error"
	default:
		jsonRet.Status = 1
		jsonRet.Msg = "other error"
	}
	data, _ := json.Marshal(jsonRet)
	ret = string(data)
	return
}
