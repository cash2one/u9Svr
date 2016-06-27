package channelPayNotify

import (
	"sort"
	"u9/tool"
)

var huaweiUrlKeys []string = []string{"requestId", "userName", "orderId",
	"amount", "result", "sign"}

//devPublicKey test
//publicKey = `MFwwDQYJKoZIhvcNAQEBBQADSwAwSAJBAIW1g+KAqqOeC1ypte8L3qTDk2nz6jUbM6o6Jg9obvivPnCAm/wZvV3jWbYWfOuO/wrFJygn/jZqf8cR1T1CQa8CAwEAAQ==`

type Huawei struct {
	Rsa
}

type huaweiTradeData struct {
	Result      string `json:"result"`
	UserName    string `json:"userName"`
	ProductName string `json:"productName"`
	PayType     string `json:"payType"`
	Amount      string `json:"amount"`
	OrderId     string `json:"orderId"`
	NotifyTime  string `json:"notifyTime"`
	RequestId   string `json:"requestId"`
	BankId      string `json:"bankId"`
	OrderTime   string `json:"orderTime"`
	TradeTime   string `json:"tradeTime"`
	AccessMode  string `json:"accessMode"`
	Spending    string `json:"spending"`
	ExtReserved string `json:"extReserved"`
	Sign        string `json:"sign"`
}

func (this *Huawei) Init(params ...interface{}) (err error) {
	if err = this.Rsa.Init(params...); err != nil {
		return
	}

	this.urlParamCheckKeys = &huaweiUrlKeys
	this.requireChannelUserId = true
	this.exChangeRatio = 100

	this.channelParamKeys["_gameId"] = ""
	this.channelParamKeys["_gameKey"] = ""
	this.channelParamKeys["_payKey"] = "HUAWEI_PAY_PUBLIC_KEY"

	this.channelTradeData = nil
	this.channelRetData = nil

	this.signMode = 0
	return
}

func (this *Huawei) ParseInputParam(params ...interface{}) (err error) {
	if err = this.Rsa.ParseInputParam(); err != nil {
		return
	}

	orderId_key := "requestId"
	channelUserId_key := "userName"
	channelOrderId_key := "orderId"
	amount_key := "amount"
	discount_key := ""
	return parseTradeData_urlParam(&this.Rsa.Base,
		orderId_key, channelUserId_key, channelOrderId_key,
		amount_key, discount_key)
}

func (this *Huawei) CheckSign(params ...interface{}) (err error) {
	excludeItems := []string{"sign"}
	sorter := tool.NewUrlValuesSorter(this.urlParams, &excludeItems)
	sort.Sort(sorter)
	this.signContent = sorter.DefaultBody()

	this.inputSign = this.urlParams.Get("sign")
	return this.Rsa.CheckSign()
}

func (this *Huawei) CheckChannelRet(params ...interface{}) (err error) {
	tradeState := this.urlParams.Get("result") == "0"
	tradeFailDesc := `urlParam(result)!="0"`
	return this.Rsa.CheckChannelRet(tradeState, tradeFailDesc)
}

func (this *Huawei) GetResult(params ...interface{}) (ret string) {
	format := `%s`
	succMsg := "SUCCESS"
	failMsg := "FAILURE"
	return this.Rsa.GetResult(format, succMsg, failMsg)
}
