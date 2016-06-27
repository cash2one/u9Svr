package channelPayNotify

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
)

var c07073UrlKeys []string = []string{"data"}

type C07073 struct {
	MD5
}

type c07073TradeData struct {
	Orderid     string `json:"orderid"`
	Gameid      string `json:"gameid"`
	Serverid    string `json:"serverid"`
	Uid         string `json:"uid"`
	Amount      string `json:"amount"`
	Time        uint   `json:"time"`
	Sign        string `json:"sign"`
	ExtendsInfo string `json:"extendsInfo"`
}

func (this *C07073) Init(params ...interface{}) (err error) {
	if err = this.MD5.Init(params...); err != nil {
		return
	}

	this.urlParamCheckKeys = &c07073UrlKeys
	this.requireChannelUserId = true
	this.exChangeRatio = 100

	this.channelParamKeys["_gameId"] = ""
	this.channelParamKeys["_gameKey"] = ""
	this.channelParamKeys["_payKey"] = "C07073_SECRET_KEY"

	this.channelTradeData = new(c07073TradeData)
	this.channelRetData = nil

	this.signHandleMethod = ""
	return
}

func (this *C07073) ParseInputParam(params ...interface{}) (err error) {
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

	this.channelTradeContent = this.urlParams.Get("data")
	if err = json.Unmarshal([]byte(this.channelTradeContent), &this.channelTradeData); err != nil {
		return
	}

	channelTradeData := this.channelTradeData.(*c07073TradeData)

	this.orderId = channelTradeData.ExtendsInfo
	this.channelUserId = channelTradeData.Uid
	this.channelOrderId = channelTradeData.Orderid

	amount := channelTradeData.Amount
	discount := ""

	if err = this.MD5.parsePayAmount(amount, discount); err != nil {
		err = nil
		return
	}
	return
}

func (this *C07073) CheckSign(params ...interface{}) (err error) {
	channelTradeData := this.channelTradeData.(*c07073TradeData)

	format := "amount=%s&gameid=%s&orderid=%s&serverid=%s&time=%d&uid=%s%s"
	this.signContent = fmt.Sprintf(format,
		channelTradeData.Amount,
		channelTradeData.Gameid,
		channelTradeData.Orderid,
		channelTradeData.Serverid,
		channelTradeData.Time,
		channelTradeData.Uid,
		this.channelParams["_payKey"])
	this.inputSign = channelTradeData.Sign
	return this.MD5.CheckSign()
}

func (this *C07073) GetResult(params ...interface{}) (ret string) {
	format := `%s`
	succMsg := "succ"
	failMsg := "fail"
	return this.MD5.GetResult(format, succMsg, failMsg)
}
