package channelPayNotify

import (
	"fmt"
	"github.com/astaxie/beego"
	"strings"
	"u9/tool"
)

var mumayiUrlKeys []string = []string{"uid", "orderID", "productPrice", "orderTime",
	"tradeSign", "tradeState"}

//木蚂蚁
type MuMaYi struct {
	Base
}

func (this *MuMaYi) Init(params ...interface{}) (err error) {
	if err = this.Base.Init(params...); err != nil {
		return
	}

	this.urlParamCheckKeys = &mumayiUrlKeys
	this.requireChannelUserId = true
	this.exChangeRatio = 100

	this.channelParamKeys["_gameId"] = ""
	this.channelParamKeys["_gameKey"] = ""
	this.channelParamKeys["_payKey"] = "MUMAYI_APPKEY"

	this.channelTradeData = nil
	this.channelRetData = nil

	return
}

func (this *MuMaYi) ParseInputParam(params ...interface{}) (err error) {
	if err = this.Base.ParseInputParam(); err != nil {
		return
	}

	orderId_key := "productDesc"
	channelUserId_key := "uid"
	channelOrderId_key := "orderID"
	amount_key := "productPrice"
	discount_key := ""
	return parseTradeData_urlParam(&this.Base,
		orderId_key, channelUserId_key, channelOrderId_key,
		amount_key, discount_key)
}

func (this *MuMaYi) CheckSign(params ...interface{}) (err error) {
	content := this.channelOrderId
	inputSign := this.urlParams.Get("tradeSign")

	var result string
	if result, err = tool.MMYSign(inputSign, this.channelParams["_payKey"], content); err == nil {
		format := "MuMaYi:%v"
		msg := fmt.Sprintf(format, err)
		beego.Error(msg)
		result = strings.TrimSpace(result)
	}

	signMethod := "MMYSign(jar)"
	format := "content:%s, inputSign:%s, result:%s"
	signMsg := fmt.Sprintf(format, content, inputSign, result)
	signState := result == "true"
	return this.Base.CheckSign(signState, signMethod, signMsg)
}

func (this *MuMaYi) CheckChannelRet(params ...interface{}) (err error) {
	tradeState := this.urlParams.Get("tradeState") == "success"
	tradeFailDesc := `urlParam(tradeState)!="success"`
	return this.Base.CheckChannelRet(tradeState, tradeFailDesc)
}

func (this *MuMaYi) GetResult(params ...interface{}) (ret string) {
	format := `%s`
	succMsg := "success"
	failMsg := "failure"
	return this.Base.GetResult(format, succMsg, failMsg)
}
