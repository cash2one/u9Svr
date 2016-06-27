package channelPayNotify

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"u9/tool"
)

var anZhiUrlKeys []string = []string{"data"}

//安智
type AnZhi struct {
	Base
}

type anZhiTradeData struct {
	PayAmount    string `json:"payAmount"`    //支付金额
	Uid          string `json:"uid"`          //用户id
	NotifyTime   int    `json:"notifyTime"`   //请求时间
	CpInfo       string `json:"cpInfo"`       //回调信息
	Memo         string `json:"memo"`         //备注
	OrderAmount  string `json:"orderAmount"`  //订单金额
	OrderAccount string `json:"orderAccount"` //订单数量
	Code         int    `json:"code"`         //订单状态
	OrderTime    string `json:"orderTime"`    //订单时间
	Msg          string `json:"msg"`          //消息
	OrderId      string `json:"orderId"`      //订单号
	RedBagMoney  string `redBagMoney`         //礼券
}

func (this *AnZhi) Init(params ...interface{}) (err error) {
	if err = this.Base.Init(params...); err != nil {
		return
	}

	this.urlParamCheckKeys = &anZhiUrlKeys
	this.requireChannelUserId = true
	this.exChangeRatio = 1

	this.channelParamKeys["_gameId"] = ""
	this.channelParamKeys["_gameKey"] = ""
	this.channelParamKeys["_payKey"] = "ANZHI_APPSECRET"

	this.channelTradeData = new(anZhiTradeData)
	this.channelRetData = nil

	return
}

func (this *AnZhi) ParseInputParam(params ...interface{}) (err error) {
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

	payKey := this.channelParams["_payKey"]
	data := this.urlParams.Get("data")
	if this.channelTradeContent, err = tool.JavaDesDecyrpt(payKey, data); err != nil {
		return
	}

	if err = json.Unmarshal([]byte(this.channelTradeContent), &this.channelTradeData); err != nil {
		return
	}

	channelTradeData := this.channelTradeData.(*anZhiTradeData)
	this.orderId = channelTradeData.CpInfo
	this.channelUserId = channelTradeData.Uid
	this.channelOrderId = channelTradeData.OrderId

	amount := channelTradeData.OrderAmount
	discount := channelTradeData.RedBagMoney

	if err = this.Base.parsePayAmount(amount, discount); err != nil {
		err = nil
		return
	}

	return
}

func (this *AnZhi) CheckChannelRet(params ...interface{}) (err error) {
	channelTradeData := this.channelTradeData.(*anZhiTradeData)
	tradeState := channelTradeData.Code == 1
	tradeFailDesc := `channelTradeData.Code!=1`
	if err = this.Base.CheckChannelRet(tradeState, tradeFailDesc); err != nil {
		return
	}
	return
}

func (this *AnZhi) CheckSign(params ...interface{}) (err error) {
	return
}

func (this *AnZhi) GetResult(params ...interface{}) (ret string) {
	format := `%s`
	succMsg := "success"
	failMsg := "failure"
	return this.Base.GetResult(format, succMsg, failMsg)
}
