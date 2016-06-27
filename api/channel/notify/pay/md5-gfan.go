package channelPayNotify

import (
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
)

var gfanUrlKeys []string = []string{"time", "sign"}

type gFanTradeData struct {
	OrderId    string `xml:"order_id"`
	Cost       int    `xml:"cost"`
	Appkey     uint64 `xml:"appkey"`
	CreateTime uint64 `xml:"create_time"`
}

type GFan struct {
	MD5
}

func (this *GFan) Init(params ...interface{}) (err error) {
	if err = this.MD5.Init(params...); err != nil {
		return
	}

	this.urlParamCheckKeys = &gfanUrlKeys
	this.requireChannelUserId = false
	this.exChangeRatio = 10 //1元 = 10机锋券

	this.channelParamKeys["_gameId"] = ""
	this.channelParamKeys["_gameKey"] = ""
	this.channelParamKeys["_payKey"] = "GFAN_UID"

	this.channelTradeData = new(gFanTradeData)
	this.channelRetData = nil

	this.signHandleMethod = ""
	return
}

func (this *GFan) ParseInputParam(params ...interface{}) (err error) {
	if err = this.MD5.ParseInputParam(); err != nil {
		return
	}

	defer func() {
		if err != nil {
			this.lastError = err_parseInputParam

			format := "ParseInputParam:err:%v, %s"
			msg := fmt.Sprintf(format, err, this.Dump())
			err = errors.New(msg)
			beego.Error(err)
		}
	}()

	this.channelTradeContent = this.body
	if err = xml.Unmarshal([]byte(this.channelTradeContent), &this.channelTradeData); err != nil {
		return
	}

	channelTradeData := this.channelTradeData.(*gFanTradeData)
	this.orderId = channelTradeData.OrderId
	this.channelUserId = ""
	this.channelOrderId = ""
	this.payAmount = channelTradeData.Cost * int(this.exChangeRatio)
	this.payDiscount = 0
	return
}

func (this *GFan) CheckSign(params ...interface{}) (err error) {
	this.signContent = this.channelParams["_payKey"] + this.urlParams.Get("time")
	this.inputSign = this.urlParams.Get("sign")
	return this.MD5.CheckSign()
}

func (this *GFan) GetResult(params ...interface{}) (ret string) {
	format := `<response><ErrorCode>%s</ErrorCode><ErrorDesc>%s</ErrorDesc></response>`
	succMsg := "1,success"
	failMsg := "0,fail"
	return this.Base.GetResult(format, succMsg, failMsg)
}
