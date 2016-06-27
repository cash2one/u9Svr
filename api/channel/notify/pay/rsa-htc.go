package channelPayNotify

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"net/url"
	"strings"
)

type HTC struct {
	Rsa
	bodyParams url.Values
}

type htcTradeData struct {
	Result_code   int    `json:"result_code"`
	Gmt_create    string `json:"gmt_create"`
	Real_amount   int    `json:"real_amount"`
	Result_msg    string `json:"result_msg"`
	Game_code     string `json:"game_code"`
	Game_order_id string `json:"game_order_id"`
	Jolo_order_id string `json:"jolo_order_id"`
	Gmt_payment   string `json:"gmt_payment"`
}

func (this *HTC) Init(params ...interface{}) (err error) {
	if err = this.Rsa.Init(params...); err != nil {
		return
	}

	this.urlParamCheckKeys = &emptyUrlKeys
	this.requireChannelUserId = false
	this.exChangeRatio = 1

	this.channelParamKeys["_gameId"] = ""
	this.channelParamKeys["_gameKey"] = ""
	this.channelParamKeys["_payKey"] = "HTC_SDK_PUBLICKEY"

	this.channelTradeData = new(htcTradeData)
	this.channelRetData = nil

	this.signMode = 0
	return
}

func (this *HTC) ParseInputParam(params ...interface{}) (err error) {
	if err = this.Rsa.ParseInputParam(); err != nil {
		return
	}

	defer func() {
		if err != nil {
			this.lastError = err_parseInputParam

			format := "ParseInputParam:err:%v, bodyParams:%+v"
			msg := fmt.Sprintf(format, err, this.bodyParams)
			err = errors.New(msg)
			beego.Error(err)
		}
	}()

	this.channelTradeContent = this.body
	if this.bodyParams, err = url.ParseQuery(this.body); err != nil {
		return
	}

	this.channelTradeContent = this.bodyParams.Get("order")
	this.channelTradeContent = strings.Replace(this.channelTradeContent, "\"{", "{", 1)
	this.channelTradeContent = strings.Replace(this.channelTradeContent, "}\"", "}", 1)

	if err = json.Unmarshal([]byte(this.channelTradeContent), &this.channelTradeData); err != nil {
		return
	}

	channelTradeData := this.channelTradeData.(*htcTradeData)

	this.orderId = channelTradeData.Game_order_id
	this.channelUserId = ""
	this.channelOrderId = channelTradeData.Jolo_order_id
	this.payAmount = channelTradeData.Real_amount * int(this.exChangeRatio)
	this.payDiscount = 0
	return
}

func (this *HTC) CheckSign(params ...interface{}) (err error) {
	this.signContent = this.channelTradeContent
	//this.bodyParams.Get("sign_type")
	this.inputSign = this.bodyParams.Get("sign")
	this.inputSign = strings.Replace(this.inputSign, "\"", "", 2)
	return this.Rsa.CheckSign()
}

func (this *HTC) CheckChannelRet(params ...interface{}) (err error) {
	channelTradeData := this.channelTradeData.(*htcTradeData)
	tradeState := channelTradeData.Result_code == 1
	tradeFailDesc := `channelTradeData.Result_code!=1`
	return this.Rsa.CheckChannelRet(tradeState, tradeFailDesc)
}

func (this *HTC) GetResult(params ...interface{}) (ret string) {
	format := `%s`
	succMsg := "success"
	failMsg := "failure"
	return this.Rsa.GetResult(format, succMsg, failMsg)
}
