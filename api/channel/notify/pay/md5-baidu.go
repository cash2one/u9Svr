package channelPayNotify

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"u9/tool"
)

type Baidu struct {
	MD5
}

type baiduTradeData struct {
	UID             string `json:"UID"`
	MerchandiseName string `json:"MerchandiseName"`
	OrderMoney      string `json:"OrderMoney"`
	StartDateTime   string `json:"StartDateTime"`
	BankDateTime    string `json:"BankDateTime"`
	OrderStatus     int    `json:"OrderStatus"`
	StatusMsg       string `json:"StatusMsg"`
	ExtInfo         string `json:"ExtInfo"`
	VoucherMoney    int    `json:"VoucherMoney"`
}

type baiduChannelRet struct {
	AppID      string `json:"AppID"`
	ResultCode string `json:"ResultCode"`
	ResultMsg  string `json:"ResultMsg"`
	Sign       string `json:"Sign"`
}

func (this *Baidu) Init(params ...interface{}) (err error) {
	if err = this.MD5.Init(params...); err != nil {
		return
	}

	this.urlParamCheckKeys = &emptyUrlKeys
	this.requireChannelUserId = true
	this.exChangeRatio = 100

	this.channelParamKeys["_gameId"] = "BAIDU_APPID"
	this.channelParamKeys["_gameKey"] = ""
	this.channelParamKeys["_payKey"] = "BAIDU_APPSECRET"

	this.channelTradeData = new(baiduTradeData)
	this.channelRetData = new(baiduChannelRet)

	this.signHandleMethod = ""
	return
}

func (this *Baidu) ParseInputParam(params ...interface{}) (err error) {
	if err = this.MD5.ParseInputParam(); err != nil {
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

	var decodeByte []byte
	content := this.urlParams.Get("Content")
	if decodeByte, err = base64.StdEncoding.DecodeString(content); err != nil {
		return
	}

	this.channelTradeContent = string(decodeByte)
	if err = json.Unmarshal(decodeByte, &this.channelTradeData); err != nil {
		return
	}

	this.orderId = this.urlParams.Get("CooperatorOrderSerial")

	channelTradeData := this.channelTradeData.(*baiduTradeData)
	this.channelUserId = channelTradeData.UID

	this.channelOrderId = this.urlParams.Get("OrderSerial")

	amount := channelTradeData.OrderMoney
	discount := ""

	if err = this.MD5.parsePayAmount(amount, discount); err != nil {
		err = nil
		return
	}
	return
}

func (this *Baidu) CheckSign(params ...interface{}) (err error) {
	this.signContent = this.urlParams.Get("AppID") +
		this.urlParams.Get("OrderSerial") +
		this.urlParams.Get("CooperatorOrderSerial") +
		this.urlParams.Get("Content") +
		this.channelParams["_payKey"]
	this.inputSign = this.urlParams.Get("Sign")
	return this.MD5.CheckSign()
}

func (this *Baidu) CheckChannelRet(params ...interface{}) (err error) {
	channelTradeData := this.channelTradeData.(*baiduTradeData)

	tradeState := this.urlParams.Get("AppID") == this.channelParams["_gameId"]
	tradeFailDesc := `urlParam(AppID)!=channelParam(_gameId)`
	if !tradeState {
		return this.MD5.CheckChannelRet(tradeState, tradeFailDesc)
	}

	tradeState = channelTradeData.OrderStatus == 1
	tradeFailDesc = `channelTradeData.OrderStatus!=1`
	return this.MD5.CheckChannelRet(tradeState, tradeFailDesc)
}

func (this *Baidu) GetResult(params ...interface{}) (ret string) {
	channelRetData := this.channelRetData.(*baiduChannelRet)
	channelRetData.AppID = this.urlParams.Get("AppID")

	if this.lastError == err_noerror {
		channelRetData.ResultCode = "1"
		channelRetData.ResultMsg = "成功"
	} else if this.lastError == err_checkSign {
		channelRetData.ResultCode = "1001"
		channelRetData.ResultMsg = "Sign无效"
	} else {
		channelRetData.ResultCode = "0"
		channelRetData.ResultMsg = "其它"
	}

	content := channelRetData.AppID + channelRetData.ResultCode + this.channelParams["_payKey"]
	channelRetData.Sign = tool.Md5([]byte(content))

	channelRetJson, _ := json.Marshal(channelRetData)
	ret = string(channelRetJson)

	this.MD5.GetResult()
	return
}
