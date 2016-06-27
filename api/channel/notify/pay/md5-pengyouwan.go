package channelPayNotify

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"net/url"
)

//朋友玩

type pengyouwanTradeData struct {
	Ver       string `json:"ver"`
	Tid       string `json:"tid"`
	Sign      string `json:"sign"`
	Gamekey   string `json:"gamekey"`
	Channel   string `json:"channel"`
	CPOrderid string `json:"cp_orderid"`
	ChOrderid string `json:"ch_orderid"`
	Amount    string `json:"amount"`
	CpParam   string `json:"cp_param"`
}

type Pengyouwan struct {
	MD5
}

func (this *Pengyouwan) Init(params ...interface{}) (err error) {
	if err = this.MD5.Init(params...); err != nil {
		return
	}

	this.urlParamCheckKeys = &emptyUrlKeys
	this.requireChannelUserId = false
	this.exChangeRatio = 100

	this.channelParamKeys["_gameId"] = ""
	this.channelParamKeys["_gameKey"] = "PYW_GAME_KEY"
	this.channelParamKeys["_payKey"] = "PYW_PAY_KEY"

	this.channelTradeData = new(pengyouwanTradeData)
	this.channelRetData = nil

	this.signHandleMethod = ""
	return
}

func (this *Pengyouwan) ParseInputParam(params ...interface{}) (err error) {
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

	this.channelTradeContent, _ = url.QueryUnescape(this.urlParams.Encode())
	this.channelTradeContent = this.channelTradeContent[0 : len(this.channelTradeContent)-1]
	if err = json.Unmarshal([]byte(this.channelTradeContent), &this.channelTradeData); err != nil {
		return
	}

	channelTradeData := this.channelTradeData.(*pengyouwanTradeData)
	this.orderId = channelTradeData.CPOrderid
	this.channelUserId = ""
	this.channelOrderId = channelTradeData.ChOrderid

	amount := channelTradeData.Amount
	discount := ""

	if err = this.MD5.parsePayAmount(amount, discount); err != nil {
		err = nil
		return
	}
	return
}

func (this *Pengyouwan) CheckSign(params ...interface{}) (err error) {
	channelTradeData := this.channelTradeData.(*pengyouwanTradeData)
	this.signContent = this.channelParams["_payKey"] +
		channelTradeData.CPOrderid +
		channelTradeData.ChOrderid +
		channelTradeData.Amount
	this.inputSign = channelTradeData.Sign
	return this.MD5.CheckSign()
}

func (this *Pengyouwan) CheckChannelRet(params ...interface{}) (err error) {
	channelTradeData := this.channelTradeData.(*pengyouwanTradeData)

	tradeState := channelTradeData.Gamekey == this.channelParams["_gameKey"]
	tradeFailDesc := `channelTradeData.Gamekey!=channelParam(_gameKey)`
	return this.MD5.CheckChannelRet(tradeState, tradeFailDesc)
}

func (this *Pengyouwan) GetResult(params ...interface{}) (ret string) {
	format := `{"ack":%s,"msg":"%s"}`
	succMsg := "200,ok"
	failMsg := "500,fail"
	return this.MD5.GetResult(format, succMsg, failMsg)
}
