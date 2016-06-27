package channelPayNotify

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"net/url"
	"u9/tool"
)

//tt

type ttTradeData struct {
	Uid        int    `json:"uid"`
	GameId     int    `json:"gameId"`
	SDKOrderId string `json:"sdkOrderId"`
	CpOrderId  string `json:"cpOrderId"`
	PayFee     string `json:"payFee"`
	PayResult  string `json:"payResult"`
	PayDate    string `json:"payDate"`
	ExInfo     string `json:"exInfo"`
}

type TT struct {
	Base
}

func (this *TT) Init(params ...interface{}) (err error) {
	if err = this.Base.Init(params...); err != nil {
		return
	}

	this.urlParamCheckKeys = &emptyUrlKeys
	this.requireChannelUserId = false
	this.exChangeRatio = 100

	this.channelParamKeys["_gameId"] = ""
	this.channelParamKeys["_gameKey"] = ""
	this.channelParamKeys["_payKey"] = "TT_SDK_PAYKEY"

	this.channelTradeData = new(ttTradeData)
	this.channelRetData = nil
	return
}

func (this *TT) ParseInputParam(params ...interface{}) (err error) {
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

	if this.channelTradeContent, err = url.QueryUnescape(this.body); err != nil {
		return
	}

	if err = json.Unmarshal([]byte(this.channelTradeContent), &this.channelTradeData); err != nil {
		return
	}

	channelTradeData := this.channelTradeData.(*ttTradeData)
	this.orderId = channelTradeData.CpOrderId
	this.channelUserId = ""
	this.channelOrderId = channelTradeData.SDKOrderId

	amount := channelTradeData.PayFee
	discount := ""

	if err = this.Base.parsePayAmount(amount, discount); err != nil {
		err = nil
		return
	}

	return
}

func (this *TT) CheckSign(params ...interface{}) (err error) {
	content := fmt.Sprintf("%s%s", this.channelTradeContent, this.channelParams["_payKey"])
	inputSign := this.ctx.Request.Header.Get("Sign")

	var sign string
	sign, err = tool.TTSign(content)
	signMethod := "TTSign(jar)"
	format := "content:%s, inputSign:%s, sign:%s"
	signMsg := fmt.Sprintf(format, content, inputSign, sign)
	signState := sign == inputSign
	return this.Base.CheckSign(signState, signMethod, signMsg)
}

func (this *TT) CheckChannelRet(params ...interface{}) (err error) {
	channelTradeData := this.channelTradeData.(*ttTradeData)
	tradeState := channelTradeData.PayResult == "1"
	tradeFailDesc := `channelTradeData.PayResult!="1"`
	return this.Base.CheckChannelRet(tradeState, tradeFailDesc)
}

func (this *TT) GetResult(params ...interface{}) (ret string) {
	format := `{"head":{"result":"%s","message":"%s"}}`
	succMsg := "0,成功"
	failMsg := "1,失败"
	return this.Base.GetResult(format, succMsg, failMsg)
}
