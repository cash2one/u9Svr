package channelPayNotify

import (
	"encoding/json"
	"fmt"
)

type m4399ChannelRet struct {
	Status    int    `json:"status"` //1：异常；2：成功；3：失败（将钱返还给用户）
	Code      string `json:"code"`
	Money     string `json:"money"`
	Gamemoney string `json:"gamemoney"`
	Msg       string `json:"msg"`
}

var m4399UrlKeys []string = []string{"orderid", "p_type", "uid", "money",
	"gamemoney", "mark", "time", "sign"}

type M4399 struct {
	MD5
}

func (this *M4399) Init(params ...interface{}) (err error) {
	if err = this.MD5.Init(params...); err != nil {
		return
	}

	this.urlParamCheckKeys = &m4399UrlKeys
	this.requireChannelUserId = true
	this.exChangeRatio = 100

	this.channelParamKeys["_gameId"] = ""
	this.channelParamKeys["_gameKey"] = ""
	this.channelParamKeys["_payKey"] = "M4399_SECRECT"

	this.channelTradeData = nil
	this.channelRetData = new(m4399ChannelRet)

	this.signHandleMethod = ""
	return
}

func (this *M4399) ParseInputParam(params ...interface{}) (err error) {
	if err = this.MD5.ParseInputParam(); err != nil {
		return
	}

	orderId_key := "mark"
	channelUserId_key := "uid"
	channelOrderId_key := "orderid"
	amount_key := "money"
	discount_key := ""
	return parseTradeData_urlParam(&this.Base,
		orderId_key, channelUserId_key, channelOrderId_key,
		amount_key, discount_key)
}

func (this *M4399) CheckSign(params ...interface{}) (err error) {
	format := `%s%s%s%s%s%s%s%s%s`
	this.signContent = fmt.Sprintf(format,
		this.channelOrderId,
		this.channelUserId,
		this.urlParams.Get("money"),
		this.urlParams.Get("gamemoney"),
		this.urlParams.Get("serverid"),
		this.channelParams["_payKey"],
		this.orderId,
		this.urlParams.Get("roleid"),
		this.urlParams.Get("time"))
	this.inputSign = this.urlParams.Get("sign")
	return this.MD5.CheckSign()
}

func (this *M4399) GetResult(params ...interface{}) (ret string) {
	channelRetData := this.channelRetData.(*m4399ChannelRet)
	channelRetData.Status = 1
	channelRetData.Money = this.urlParams.Get("money")
	channelRetData.Gamemoney = this.urlParams.Get("gamemoney")

	switch this.lastError {
	case err_noerror:
		channelRetData.Status = 2
		channelRetData.Msg = "success"
	case err_checkSign:
		channelRetData.Msg = "sign_error"
	case err_prepareLoginRequest:
		fallthrough
	case err_channelUserIsNotExist:
		channelRetData.Msg = "user_not_exist"
	case err_handleOrder:
		channelRetData.Msg = "orderid_exist"
	case err_payAmountError:
		channelRetData.Msg = "money_error"
	default:
		channelRetData.Msg = "other error"
	}
	jsonByte, _ := json.Marshal(channelRetData)
	ret = string(jsonByte)

	this.MD5.GetResult()
	return
}
