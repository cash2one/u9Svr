package channelPayNotify

import (
	"sort"
	"u9/tool"
)

var qihoo360UrlKeys []string = []string{"app_key", "product_id", "amount", "app_uid",
	"user_id", "order_id", "gateway_flag", "sign_type", "sign_return", "sign"}

type Qihoo360 struct {
	MD5
}

type qihoo360TradeData struct {
	app_key      string `json:"app_key"`
	product_id   int    `json:"product_id"`
	amount       int    `json:"amount"`
	app_uid      string `json:"app_uid"`
	app_ext1     int    `json:"app_ext1"`
	app_ext2     string `json:"app_ext2"`
	user_id      int    `json:"user_id"`
	order_id     int    `json:"order_id"`
	gateway_flag string `json:"gateway_flag"`
	sign_type    string `json:"sign_type"`
	app_order_id string `json:"app_order_id"`
	sign_return  string `json:"sign_return"`
	sign         string `json:"sign"`
}

func (this *Qihoo360) Init(params ...interface{}) (err error) {
	if err = this.MD5.Init(params...); err != nil {
		return
	}

	this.urlParamCheckKeys = &qihoo360UrlKeys
	this.requireChannelUserId = true
	this.exChangeRatio = 1

	this.channelParamKeys["_gameId"] = ""
	this.channelParamKeys["_gameKey"] = ""
	this.channelParamKeys["_payKey"] = "QIHOO360_SECRET"

	this.channelTradeData = nil
	this.channelRetData = nil

	this.signHandleMethod = ""
	return
}

func (this *Qihoo360) ParseInputParam(params ...interface{}) (err error) {
	if err = this.MD5.ParseInputParam(); err != nil {
		return
	}

	orderId_key := "app_order_id"
	channelUserId_key := "user_id"
	channelOrderId_key := "order_id"
	amount_key := "amount"
	discount_key := ""
	return parseTradeData_urlParam(&this.Base,
		orderId_key, channelUserId_key, channelOrderId_key,
		amount_key, discount_key)
}

func (this *Qihoo360) CheckSign(params ...interface{}) (err error) {
	excludeItems := []string{"sign_return", "sign"}
	sorter := tool.NewUrlValuesSorter(this.urlParams, &excludeItems)
	sort.Sort(sorter)
	this.signContent = sorter.FormatBody("v", "#") + "#" + this.channelParams["_payKey"]

	this.inputSign = this.urlParams.Get("sign")
	return this.MD5.CheckSign()
}

func (this *Qihoo360) GetResult(params ...interface{}) (ret string) {
	format := `%s`
	succMsg := "ok"
	failMsg := "fail"
	return this.MD5.GetResult(format, succMsg, failMsg)
}
