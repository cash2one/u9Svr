package channelPayNotify

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"strings"
	"u9/tool"
)

var coolpadUrlKeys []string = []string{"transdata", "sign"}

type coolPadTradeData struct {
	Transtype int     `json:"transtype"`
	Cporderid string  `json:"cporderid"`
	Transid   string  `json:"transid"`
	Appuserid string  `json:"appuserid"`
	Appid     string  `json:"appid"`
	Waresid   int     `json:"waresid"`
	Feetype   int     `json:"feetype"`
	Money     float64 `json:"money`
	Currency  string  `json:"currency"`
	Result    int     `json:"result"`
	Transtime string  `json:"transtime"`
	Cpprivate string  `json:"cpprivate"`
	PayType   int     `json:"paytype"`
	Sign      string  `json:"Sign"`
	Signtype  string  `json:"Signtype"`
	Code      string  `json:"code"`
	ErrMsg    string  `json:"errmsg"`
}

type CoolPad struct {
	Base
}

func (this *CoolPad) Init(params ...interface{}) (err error) {
	if err = this.Base.Init(params...); err != nil {
		return
	}

	this.urlParamCheckKeys = &coolpadUrlKeys
	this.requireChannelUserId = true
	this.exChangeRatio = 100

	this.channelParamKeys["_gameId"] = ""
	this.channelParamKeys["_gameKey"] = ""
	this.channelParamKeys["_payKey"] = "COOLPAD_PUBLICKEY"

	this.channelTradeData = new(coolPadTradeData)
	this.channelRetData = nil

	return
}

func (this *CoolPad) ParseInputParam(params ...interface{}) (err error) {
	if err = this.Base.ParseInputParam(); err != nil {
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

	this.channelTradeContent = this.urlParams.Get("transdata")
	if err = json.Unmarshal([]byte(this.channelTradeContent), &this.channelTradeData); err != nil {
		return
	}

	channelTradeData := this.channelTradeData.(*coolPadTradeData)

	this.orderId = channelTradeData.Cporderid
	this.channelUserId = channelTradeData.Appuserid
	this.channelOrderId = channelTradeData.Transid
	this.payAmount = int(channelTradeData.Money * this.exChangeRatio)
	this.payDiscount = 0
	return
}

func (this *CoolPad) CheckSign(params ...interface{}) (err error) {
	inputSign := this.urlParams.Get("sign")
	payKey := this.channelParams["_payKey"]

	var result string
	if result, err = tool.IapppayVerify(this.channelTradeContent, inputSign, payKey); err != nil {
		format := "IapppayVerify:%v"
		msg := fmt.Sprintf(format, err)
		beego.Error(msg)
	} else {
		result = strings.TrimSpace(result)
	}

	signMethod := "CoolPadSign7.0.0(jar)"
	format := "channelTradeContent:%s, inputSign:%s, result:%s"
	signMsg := fmt.Sprintf(format, this.channelTradeContent, inputSign, result)
	signState := result == "0"
	return this.Base.CheckSign(signState, signMethod, signMsg)
}

func (this *CoolPad) CheckChannelRet(params ...interface{}) (err error) {
	channelTradeData := this.channelTradeData.(*coolPadTradeData)
	tradeState := channelTradeData.Result == 0
	tradeFailDesc := `channelTradeData.Result!=0`
	return this.Base.CheckChannelRet(tradeState, tradeFailDesc)
}

func (this *CoolPad) GetResult(params ...interface{}) (ret string) {
	format := `%s`
	succMsg := "SUCCESS"
	failMsg := "FAILURE"
	return this.Base.GetResult(format, succMsg, failMsg)
}
