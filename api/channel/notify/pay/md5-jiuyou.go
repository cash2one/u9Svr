package channelPayNotify

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
)

type Jiuyou struct {
	MD5
}

type jiuyouTradeData struct {
	Ver  string `json:"ver"`
	Data struct {
		OrderId      string `json:"orderId"`
		GameId       string `json:"gameId"`
		AccountId    string `json:"accountId"`
		Creator      string `json:"creator"`
		PayWay       string `json:"payWay"`
		Amount       string `json:"amount"`
		CallbackInfo string `json:"callbackInfo"`
		OrderStatus  string `json:"orderStatus"`
		FailedDesc   string `json:"failedDesc"`
		CpOrderId    string `json:"cpOrderId"`
	} `json:"data"`
	Sign string `json:"sign"`
}

func (this *Jiuyou) Init(params ...interface{}) (err error) {
	if err = this.MD5.Init(params...); err != nil {
		return
	}

	this.urlParamCheckKeys = &emptyUrlKeys
	this.requireChannelUserId = true
	this.exChangeRatio = 100

	this.channelParamKeys["_gameId"] = ""
	this.channelParamKeys["_gameKey"] = ""
	this.channelParamKeys["_payKey"] = "UC_APPKEY"

	this.channelTradeData = new(jiuyouTradeData)
	this.channelRetData = nil

	this.signHandleMethod = ""
	return
}

func (this *Jiuyou) ParseInputParam(params ...interface{}) (err error) {
	if err = this.MD5.ParseInputParam(); err != nil {
		return
	}

	defer func() {
		if err != nil {
			this.lastError = err_parseInputParam

			format := "ParseInputParam:err:%v"
			msg := fmt.Sprintf(format, err)
			err = errors.New(msg)
			beego.Error(err)
		}
	}()

	this.channelTradeContent = this.body
	if err = json.Unmarshal([]byte(this.channelTradeContent), &this.channelTradeData); err != nil {
		return
	}

	channelTradeData := this.channelTradeData.(*jiuyouTradeData)
	this.orderId = channelTradeData.Data.CpOrderId
	this.channelUserId = channelTradeData.Data.AccountId
	this.channelOrderId = channelTradeData.Data.OrderId

	amount := channelTradeData.Data.Amount
	discount := ""

	if err = this.MD5.parsePayAmount(amount, discount); err != nil {
		err = nil
		return
	}

	return
}

func (this *Jiuyou) CheckSign(params ...interface{}) (err error) {
	channelTradeData := this.channelTradeData.(*jiuyouTradeData)
	format := `accountId=%samount=%scallbackInfo=%scpOrderId=%screator=%s` +
		`failedDesc=%sgameId=%sorderId=%sorderStatus=%spayWay=%s%s`
	this.signContent = fmt.Sprintf(format,
		channelTradeData.Data.AccountId,
		channelTradeData.Data.Amount,
		channelTradeData.Data.CallbackInfo,
		channelTradeData.Data.CpOrderId,
		channelTradeData.Data.Creator,
		channelTradeData.Data.FailedDesc,
		channelTradeData.Data.GameId,
		channelTradeData.Data.OrderId,
		channelTradeData.Data.OrderStatus,
		channelTradeData.Data.PayWay,
		this.channelParams["_payKey"])
	this.inputSign = channelTradeData.Sign
	return this.MD5.CheckSign()
}

func (this *Jiuyou) CheckChannelRet(params ...interface{}) (err error) {
	channelTradeData := this.channelTradeData.(*jiuyouTradeData)
	tradeState := channelTradeData.Data.OrderStatus == "S"
	tradeFailDesc := `channelTradeData.Data.OrderStatus!="S"`
	return this.MD5.CheckChannelRet(tradeState, tradeFailDesc)
}

func (this *Jiuyou) GetResult(params ...interface{}) (ret string) {
	format := `%s`
	succMsg := "SUCCESS"
	failMsg := "FAILURE"
	return this.MD5.GetResult(format, succMsg, failMsg)
}
