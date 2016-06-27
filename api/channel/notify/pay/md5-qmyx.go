package channelPayNotify

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"strconv"
)

type Qmyx struct {
	MD5
}

type qmyxTradeData struct {
	ChannelOrderId string `json:"pay_order_code"`
	GameId         string `json:"game_id"`
	ServerId       string `json:"server_id"`
	UserKey        string `json:"user_key"`
	Amount         int    `json:"amount"`
	Ext            string `json:"ext"`
}

func (this *Qmyx) Init(params ...interface{}) (err error) {
	if err = this.MD5.Init(params...); err != nil {
		return
	}

	this.urlParamCheckKeys = &emptyUrlKeys
	this.requireChannelUserId = true
	this.exChangeRatio = 1

	this.channelParamKeys["_gameId"] = "GAME_ID"
	this.channelParamKeys["_gameKey"] = ""
	this.channelParamKeys["_payKey"] = "GAME_KEY"

	this.channelTradeData = new(qmyxTradeData)
	this.channelRetData = nil

	this.signHandleMethod = ""
	return
}

func (this *Qmyx) ParseInputParam(params ...interface{}) (err error) {
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

	this.channelTradeContent = this.urlParams.Get("data")
	if err = json.Unmarshal([]byte(this.channelTradeContent), &this.channelTradeData); err != nil {
		return err
	}

	data := this.channelTradeData.(*qmyxTradeData)
	this.orderId = data.Ext
	this.channelOrderId = data.ChannelOrderId
	this.channelUserId = data.UserKey
	this.payAmount = data.Amount * int(this.exChangeRatio)

	return
}

func (this *Qmyx) CheckSign(params ...interface{}) (err error) {
	data := this.channelTradeData.(*qmyxTradeData)
	this.signContent = data.ChannelOrderId +
		data.GameId +
		data.ServerId +
		data.UserKey +
		strconv.Itoa(data.Amount) +
		data.Ext +
		this.channelParams["_payKey"]
	this.inputSign = this.urlParams.Get("sign")
	return this.MD5.CheckSign()
}

func (this *Qmyx) CheckChannelRet(params ...interface{}) (err error) {
	channelTradeData := this.channelTradeData.(*qmyxTradeData)
	tradeState := channelTradeData.GameId == this.channelParams["_gameId"]
	tradeFailDesc := `channelTradeData.GameId!=channelParam(_gameId)`
	return this.MD5.CheckChannelRet(tradeState, tradeFailDesc)
}

func (this *Qmyx) GetResult(params ...interface{}) (ret string) {
	format := `%s`
	succMsg := "success"
	failMsg := "failure"
	return this.MD5.GetResult(format, succMsg, failMsg)
}
