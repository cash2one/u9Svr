package channelPayNotify

import (
	"fmt"
)

var zhuoyiUrlKeys []string = []string{"Recharge_Id", "App_Id", "Uin", "Urecharge_Id",
	"Recharge_Money", "Pay_Status", "Create_Time", "Sign"}

type ZhuoYi struct {
	MD5
}

func (this *ZhuoYi) Init(params ...interface{}) (err error) {
	if err = this.MD5.Init(params...); err != nil {
		return
	}

	this.urlParamCheckKeys = &zhuoyiUrlKeys
	this.requireChannelUserId = true
	this.exChangeRatio = 100

	this.channelParamKeys["_gameId"] = ""
	this.channelParamKeys["_gameKey"] = ""
	this.channelParamKeys["_payKey"] = "zy_app_secret"

	this.channelTradeData = nil
	this.channelRetData = nil

	this.signHandleMethod = ""
	return
}

func (this *ZhuoYi) ParseInputParam(params ...interface{}) (err error) {
	if err = this.MD5.ParseInputParam(); err != nil {
		return
	}

	orderId_key := "Urecharge_Id"
	channelUserId_key := "Uin"
	channelOrderId_key := "Recharge_Id"
	amount_key := "Recharge_Money"
	discount_key := ""
	return parseTradeData_urlParam(&this.Base,
		orderId_key, channelUserId_key, channelOrderId_key,
		amount_key, discount_key)
}

func (this *ZhuoYi) CheckSign(params ...interface{}) (err error) {
	format := `App_Id=%s&Create_Time=%s&Extra=%s&Pay_Status=%s&Recharge_Gold_Count=%s` +
		`&Recharge_Id=%s&Recharge_Money=%s&Uin=%s&Urecharge_Id=%s%s`
	this.signContent = fmt.Sprintf(format,
		this.urlParams.Get("App_Id"),
		this.urlParams.Get("Create_Time"),
		this.urlParams.Get("Extra"),
		this.urlParams.Get("Pay_Status"),
		this.urlParams.Get("Recharge_Gold_Count"),
		this.urlParams.Get("Recharge_Id"),
		this.urlParams.Get("Recharge_Money"),
		this.urlParams.Get("Uin"),
		this.urlParams.Get("Urecharge_Id"),
		this.channelParams["_payKey"])
	this.inputSign = this.urlParams.Get("Sign")
	return this.MD5.CheckSign()
}

func (this *ZhuoYi) CheckChannelRet(params ...interface{}) (err error) {
	tradeState := this.urlParams.Get("Pay_Status") == "1"
	tradeFailDesc := `urlParam(Pay_Status)!="1"`
	return this.MD5.CheckChannelRet(tradeState, tradeFailDesc)
}

func (this *ZhuoYi) GetResult(params ...interface{}) (ret string) {
	format := `%s`
	succMsg := "success"
	failMsg := "failure"
	return this.MD5.GetResult(format, succMsg, failMsg)
}
