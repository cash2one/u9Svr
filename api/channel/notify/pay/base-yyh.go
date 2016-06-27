package channelPayNotify

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"strings"
	"u9/tool"
)

var yyhUrlKeys []string = []string{"transdata", "sign"}

type yyhTradeData_620 struct {
	Exorderno string `json:"exorderno"`
	Transid   string `json:"transid"`
	Appid     string `json:"appid"`
	Waresid   int    `json:"waresid"`
	Feetype   int    `json:"feetype"`
	Money     int    `json:"money"`
	Count     int    `json:"count"`
	Result    int    `json:"result`
	Transtype int    `json:"transtype"`
	Transtime string `json:"transtime"`
	Cpprivate string `json:"cpprivate"`
	PayType   string `json:"paytype"`
}

type yyhTradeData_700 struct {
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

//应用汇
type YYH struct {
	Base
}

func (this *YYH) Init(params ...interface{}) (err error) {
	if err = this.Base.Init(params...); err != nil {
		return
	}

	var channelParamErr error
	if this.channelParams["_version"], channelParamErr = this.getChannelParam(ChannelApiVerKeyName); err != nil {
		format := "init: err:%v"
		msg := fmt.Sprintf(format, channelParamErr)
		beego.Warn(msg)
	}

	channelApiVersion := this.channelParams["_version"]
	switch channelApiVersion {
	case "6.2.2":
		this.urlParamCheckKeys = &yyhUrlKeys
		this.requireChannelUserId = false
		this.exChangeRatio = 1

		this.channelTradeData = new(yyhTradeData_620)
		this.channelParamKeys["_gameId"] = ""
		this.channelParamKeys["_gameKey"] = ""
		this.channelParamKeys["_payKey"] = "YYH_APPKEY"
	case "7.0.0":
		fallthrough
	default:
		this.urlParamCheckKeys = &yyhUrlKeys
		this.requireChannelUserId = true
		this.exChangeRatio = 100

		this.channelTradeData = new(yyhTradeData_700)
		this.channelParamKeys["_gameId"] = ""
		this.channelParamKeys["_gameKey"] = ""
		this.channelParamKeys["_payKey"] = "YYH_PUBLICKEY"

	}

	this.channelRetData = nil

	return
}

func (this *YYH) ParseInputParam(params ...interface{}) (err error) {
	if err = this.Base.ParseInputParam(); err != nil {
		return
	}

	channelApiVersion := this.channelParams["_version"]
	defer func() {
		if err != nil {
			this.lastError = err_parseInputParam

			format := "ParseInputParam: err:%v, channelApiVersion:%s"
			msg := fmt.Sprintf(format, err, channelApiVersion)
			err = errors.New(msg)
			beego.Error(err)
		}
	}()

	this.channelTradeContent = this.urlParams.Get("transdata")
	if err = json.Unmarshal([]byte(this.channelTradeContent), &this.channelTradeData); err != nil {
		return
	}

	switch channelApiVersion {
	case "6.2.2":
		channelTradeData := this.channelTradeData.(*yyhTradeData_620)

		this.orderId = channelTradeData.Exorderno
		this.channelUserId = ""
		this.channelOrderId = channelTradeData.Transid
		this.payAmount = channelTradeData.Money * int(this.exChangeRatio)
		this.payDiscount = 0
	case "7.0.0":
		fallthrough
	default:
		channelTradeData := this.channelTradeData.(*yyhTradeData_700)

		this.orderId = channelTradeData.Cporderid
		this.channelUserId = channelTradeData.Appuserid
		this.channelOrderId = channelTradeData.Transid
		this.payAmount = int(channelTradeData.Money * this.exChangeRatio)
		this.payDiscount = 0
	}
	return
}

func (this *YYH) CheckSign(params ...interface{}) (err error) {
	channelApiVersion := this.channelParams["_version"]
	inputSign := this.urlParams.Get("sign")
	payKey := this.channelParams["_payKey"]

	switch channelApiVersion {
	case "6.2.2":
		sign := tool.Md5([]byte(this.channelTradeContent))

		var result string
		if result, err = tool.YYHSign(sign, inputSign, payKey); err != nil {
			format := "checkSign: YYHSign:%v"
			msg := fmt.Sprintf(format, err)
			beego.Error(msg)
		} else {
			result = strings.TrimSpace(result)
		}

		signMethod := "YYHSign6.2.2(jar)"
		format := "channelTradeContent:%s, inputSign:%s, result:%s"
		signMsg := fmt.Sprintf(format, this.channelTradeContent, inputSign, result)
		signState := result == "0"
		return this.Base.CheckSign(signState, signMethod, signMsg)
	case "7.0.0":
		fallthrough
	default:

		inputSign := this.urlParams.Get("sign")
		payKey := this.channelParams["_payKey"]

		var result string
		if result, err = tool.IapppayVerify(this.channelTradeContent, inputSign, payKey); err != nil {
			format := "checkSign: IapppayVerify(jar):%v"
			msg := fmt.Sprintf(format, err)
			beego.Error(msg)
		} else {
			result = strings.TrimSpace(result)
		}

		signMethod := "IapppaySign(jar)"
		format := "channelTradeContent:%s, inputSign:%s, result:%s"
		signMsg := fmt.Sprintf(format, this.channelTradeContent, inputSign, result)
		signState := result == "0"
		return this.Base.CheckSign(signState, signMethod, signMsg)
	}
	return
}

func (this *YYH) CheckChannelRet(params ...interface{}) (err error) {
	channelApiVersion := this.channelParams["_version"]
	switch channelApiVersion {
	case "6.2.2":
		channelTradeData := this.channelTradeData.(*yyhTradeData_620)

		tradeState := channelTradeData.Result == 0
		tradeFailDesc := `channelTradeData.Result!=0`
		return this.Base.CheckChannelRet(tradeState, tradeFailDesc)
	case "7.0.0":
		fallthrough
	default:
		channelTradeData := this.channelTradeData.(*yyhTradeData_700)

		tradeState := channelTradeData.Result == 0
		tradeFailDesc := `channelTradeData.Result!=0`
		return this.Base.CheckChannelRet(tradeState, tradeFailDesc)
	}

	return
}

func (this *YYH) GetResult(params ...interface{}) (ret string) {
	format := `%s`
	succMsg := ""
	failMsg := ""
	channelApiVersion := this.channelParams["_version"]
	switch channelApiVersion {
	case "6.2.2":
		succMsg = "true"
		failMsg = "false"
	case "7.0.0":
		fallthrough
	default:
		succMsg = "SUCCESS"
		failMsg = "FAILURE"
	}
	return this.Base.GetResult(format, succMsg, failMsg)
}
