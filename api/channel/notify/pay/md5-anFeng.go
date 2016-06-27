package channelPayNotify

import (
	"sort"
	"u9/tool"
)

type AnFeng struct {
	MD5
}

type anFengResData struct {
	Uid        string `json:"uid"`
	Ucid       string `json:"ucid"`
	Body       string `json:"body"`
	Fee        string `json:"fee"`
	Subject    string `json:"subject"`
	Vid        string `json:"vid"`
	Sn         string `json:"sn"`
	Vorderid   string `json:"vorderid"`
	CreateTime string `json:"createTime"`
	Sign       string `json:"sign"`
}

func (this *AnFeng) Init(params ...interface{}) (err error) {
	if err = this.MD5.Init(params...); err != nil {
		return
	}

	this.urlParamCheckKeys = &emptyUrlKeys
	this.requireChannelUserId = true
	this.exChangeRatio = 100

	this.channelParamKeys["_gameId"] = "productid"
	this.channelParamKeys["_gameKey"] = ""
	this.channelParamKeys["_payKey"] = "ANFENG_PAYKEY"

	this.channelTradeData = nil
	this.channelRetData = nil

	this.signHandleMethod = ""
	return
}

func (this *AnFeng) ParseInputParam(params ...interface{}) (err error) {
	if err = this.MD5.ParseInputParam(); err != nil {
		return
	}

	orderId_key := "vorderid"
	channelUserId_key := "ucid"
	channelOrderId_key := "sn"
	amount_key := "fee"
	discount_key := ""
	return parseTradeData_urlParam(&this.MD5.Base,
		orderId_key, channelUserId_key, channelOrderId_key,
		amount_key, discount_key)
}

func (this *AnFeng) CheckSign(params ...interface{}) (err error) {
	excludeItems := []string{"sign"}
	sorter := tool.NewUrlValuesSorter(this.urlParams, &excludeItems)
	sort.Sort(sorter)

	this.signContent = sorter.DefaultBody() + "&signKey=" + this.channelParams["_payKey"]
	this.inputSign = this.urlParams.Get("sign")

	return this.MD5.CheckSign()
}

func (this *AnFeng) CheckChannelRet(params ...interface{}) (err error) {
	tradeState := this.urlParams.Get("vid") == this.channelParams["_gameId"]
	tradeFailDesc := `urlParam(vid)!=channelParam(_gameId)`
	return this.MD5.CheckChannelRet(tradeState, tradeFailDesc)
}

func (this *AnFeng) GetResult(params ...interface{}) (ret string) {
	format := `%s`
	succMsg := "SUCCESS"
	failMsg := "FAILURE"
	return this.MD5.GetResult(format, succMsg, failMsg)
}
