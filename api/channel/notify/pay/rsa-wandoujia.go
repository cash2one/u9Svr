package channelPayNotify

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"strconv"
)

var wandoujiaUrlKeys []string = []string{"content", "signType", "sign"}

const wandoujiaRsaPublicKey = `MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCd95FnJFhPinpNiE/h4VA6bU1rzRa5+a25BxsnFX8TzquWxqDCoe4xG6QKXMXuKvV57tTRpzRo2jeto40eHKClzEgjx9lTYVb2RFHHFWio/YGTfnqIPTVpi7d7uHY+0FZ0lYL5LlW4E2+CQMxFOPRwfqGzMjs1SDlH7lVrLEVy6QIDAQAB`

type wandoujiaTradeData struct {
	TimeStamp  uint64 `json:"timeStamp"`
	OrderId    uint64 `json:"orderId"`
	Money      int    `json:"money"`
	ChargeType string `json:"chargeType"`
	AppKeyId   uint64 `json:"appKeyId"`
	BuyerId    uint64 `json:"buyerId"`
	OutTradeNo string `json:"out_trade_no"`
	CardNo     uint64 `json:"cardNo"`
}

type Wandoujia struct {
	Rsa
}

func (this *Wandoujia) Init(params ...interface{}) (err error) {
	if err = this.Rsa.Init(params...); err != nil {
		return
	}

	this.urlParamCheckKeys = &wandoujiaUrlKeys
	this.requireChannelUserId = true
	this.exChangeRatio = 1

	this.channelParamKeys["_gameId"] = ""
	this.channelParamKeys["_gameKey"] = "WANDOUJIA_SECRETKEY"
	this.channelParamKeys["_payKey"] = ""

	this.channelParams["_payKey"] = wandoujiaRsaPublicKey

	this.channelTradeData = new(wandoujiaTradeData)
	this.channelRetData = nil

	this.signMode = 0
	return
}

func (this *Wandoujia) ParseInputParam(params ...interface{}) (err error) {
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

	this.channelTradeContent = this.urlParams.Get("content")
	if err = json.Unmarshal([]byte(this.channelTradeContent), &this.channelTradeData); err != nil {
		return
	}

	channelTradeData := this.channelTradeData.(*wandoujiaTradeData)
	this.orderId = channelTradeData.OutTradeNo
	this.channelUserId = strconv.FormatUint(channelTradeData.BuyerId, 10)
	this.channelOrderId = strconv.FormatUint(channelTradeData.OrderId, 10)
	this.payAmount = channelTradeData.Money
	this.payDiscount = 0
	return
}

func (this *Wandoujia) CheckSign(params ...interface{}) (err error) {
	this.signContent = this.channelTradeContent
	this.inputSign = this.urlParams.Get("sign")
	return this.Rsa.CheckSign()
}

func (this *Wandoujia) GetResult(params ...interface{}) (ret string) {
	format := `%s`
	succMsg := "success"
	failMsg := "fail"
	return this.Rsa.GetResult(format, succMsg, failMsg)
}
