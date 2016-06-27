package channelPayNotify

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	//"net/url"
)

type Lenovo struct {
	Rsa
}

type lenovoTradeData struct {
	TransData struct {
		ExOrderNo string `json:"exorderno"`
		TransId   string `json:"transid"`
		AppId     string `json:"appid"`
		WareSID   int    `json:"waresid"`
		FeeType   int    `json:"feetype"`
		Money     int    `json:"money"`
		Count     int    `json:"count"`
		Result    int    `json:"result"`
		TransType int    `json:"transtype"`
		TransTime string `json:"transtime"`
		CpPrivate string `json:"cpprivate"`
		PayType   int    `json:"paytype"`
	} `json:"transdata"`
	Sign string `json:"sign"`
}

func (this *Lenovo) Init(params ...interface{}) (err error) {
	if err = this.Rsa.Init(params...); err != nil {
		return
	}

	this.urlParamCheckKeys = &emptyUrlKeys
	this.requireChannelUserId = false
	this.exChangeRatio = 1

	this.channelParamKeys["_gameId"] = "lenovo.open.appid"
	this.channelParamKeys["_gameKey"] = ""
	this.channelParamKeys["_payKey"] = "lenovo.open.appkey"

	this.channelTradeData = new(lenovoTradeData)
	this.channelRetData = nil

	this.signMode = 1
	return
}

func (this *Lenovo) ParseInputParam(params ...interface{}) (err error) {
	if err = this.Rsa.ParseInputParam(); err != nil {
		return
	}

	defer func() {
		if err != nil {
			this.lastError = err_parseInputParam

			format := "ParseInputParam: err:%v"
			msg := fmt.Sprintf(format, err)
			err = errors.New(msg)
			beego.Error(err)
		}
	}()

	channelTradeData := this.channelTradeData.(*lenovoTradeData)
	channelTradeData.Sign = this.ctx.Request.FormValue("sign")

	this.channelTradeContent = this.ctx.Request.FormValue("transdata")
	if err = json.Unmarshal([]byte(this.channelTradeContent), &channelTradeData.TransData); err != nil {
		return
	}

	this.orderId = channelTradeData.TransData.ExOrderNo
	this.channelUserId = ""
	this.channelOrderId = channelTradeData.TransData.TransId
	this.payAmount = channelTradeData.TransData.Money * int(this.exChangeRatio)
	this.payDiscount = 0
	return
}

func (this *Lenovo) CheckSign(params ...interface{}) (err error) {
	channelTradeData := this.channelTradeData.(*lenovoTradeData)
	this.signContent = this.channelTradeContent
	this.inputSign = channelTradeData.Sign
	return this.Rsa.CheckSign()
}

func (this *Lenovo) CheckChannelRet(params ...interface{}) (err error) {
	channelTradeData := this.channelTradeData.(*lenovoTradeData)

	tradeState := channelTradeData.TransData.AppId == this.channelParams["_gameId"]
	tradeFailDesc := `channelTradeData.TransData.AppId!=channelParam(_gameId)`
	if !tradeState {
		return this.Rsa.CheckChannelRet(tradeState, tradeFailDesc)
	}

	tradeState = channelTradeData.TransData.Result == 0
	tradeFailDesc = `channelTradeData.TransData.Result!=0`
	return this.Rsa.CheckChannelRet(tradeState, tradeFailDesc)
}

func (this *Lenovo) GetResult(params ...interface{}) (ret string) {
	format := `%s`
	succMsg := "SUCCESS"
	failMsg := "FAILURE"
	return this.Rsa.GetResult(format, succMsg, failMsg)
}
