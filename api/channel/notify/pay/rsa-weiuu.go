package channelPayNotify

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"regexp"
	"strconv"
	"strings"
	"u9/tool"
)

type WeiUU struct {
	Base
}

type weiuuTradeData struct {
	TransData struct {
		Extension string `json:"extension"`
		UserID    int    `json:"userID"`
		OrderID   int    `json:"orderID"`
		GameID    int    `json:"gameID"`
		ChannelID int    `json:"channelID"`
		Money     int    `json:"money"`
		ServerID  string `json:"serverID"`
		ProductID string `json:"productID"`
		Currency  string `json:"currency"`
	} `json:"data"`
	State int    `json:state`
	Sign  string `json:"sign"`
}

func (this *WeiUU) Init(params ...interface{}) (err error) {
	if err = this.Base.Init(params...); err != nil {
		return
	}

	this.urlParamCheckKeys = &emptyUrlKeys
	this.requireChannelUserId = false
	this.exChangeRatio = 1

	this.channelParamKeys["_gameId"] = ""
	this.channelParamKeys["_gameKey"] = ""
	this.channelParamKeys["_payKey"] = "WEIUU_PUBLICKEY"

	this.channelTradeData = new(weiuuTradeData)
	this.channelRetData = nil

	return
}

func (this *WeiUU) ParseInputParam(params ...interface{}) (err error) {
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

	this.channelTradeContent = this.body
	if err = json.Unmarshal([]byte(this.channelTradeContent), &this.channelTradeData); err != nil {
		return
	}
	channelTradeData := this.channelTradeData.(*weiuuTradeData)
	this.orderId = channelTradeData.TransData.Extension
	this.channelUserId = strconv.Itoa(channelTradeData.TransData.UserID)
	this.channelOrderId = strconv.Itoa(channelTradeData.TransData.OrderID)
	this.payAmount = channelTradeData.TransData.Money * int(this.exChangeRatio)
	this.payDiscount = 0
	return
}

func (this *WeiUU) CheckSign(params ...interface{}) (err error) {
	channelTradeData := this.channelTradeData.(*weiuuTradeData)
	inputSign := channelTradeData.Sign
	payKey := this.channelParams["_payKey"]

	reg, _ := regexp.Compile(`"data":({[^\}]*})`)
	content := reg.FindStringSubmatch(this.channelTradeContent)[1]

	var result string
	if result, err = tool.IapppayVerify(content, inputSign, payKey); err != nil {
		format := "IapppayVerify:%v"
		msg := fmt.Sprintf(format, err)
		beego.Error(msg)
	} else {
		result = strings.TrimSpace(result)
	}

	signMethod := "CoolPadSign7.0.0(jar)"
	format := "content:%s, inputSign:%s, result:%s"
	signMsg := fmt.Sprintf(format, content, inputSign, result)
	signState := result == "0"
	return this.Base.CheckSign(signState, signMethod, signMsg)
}

func (this *WeiUU) CheckChannelRet(params ...interface{}) (err error) {
	channelTradeData := this.channelTradeData.(*weiuuTradeData)
	tradeState := channelTradeData.State == 1
	tradeFailDesc := `cchannelTradeData.State!=1`
	return this.Base.CheckChannelRet(tradeState, tradeFailDesc)
}

func (this *WeiUU) GetResult(params ...interface{}) (ret string) {
	format := `%s`
	succMsg := "SUCCESS"
	failMsg := "FAIL"
	return this.Base.GetResult(format, succMsg, failMsg)
}
