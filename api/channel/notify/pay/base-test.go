package channelPayNotify

var testUrlKeys []string = []string{"money", "order", "mid", "ext", "result"}

type Test struct {
	Base
}

func (this *Test) Init(params ...interface{}) (err error) {
	if err = this.Base.Init(params...); err != nil {
		if this.lastError != err_initChannelParam {
			return
		} else {
			this.lastError = err_noerror
			err = nil
		}
	}

	this.urlParamCheckKeys = &testUrlKeys
	this.requireChannelUserId = true
	this.exChangeRatio = 100

	this.channelParamKeys["_gameId"] = ""
	this.channelParamKeys["_gameKey"] = ""
	this.channelParamKeys["_payKey"] = ""

	this.channelTradeData = nil
	this.channelRetData = nil

	return
}

func (this *Test) ParseInputParam(params ...interface{}) (err error) {
	if err = this.Base.ParseInputParam(); err != nil {
		return
	}

	orderId_key := "ext"
	channelUserId_key := "mid"
	channelOrderId_key := "order"
	amount_key := "money"
	discount_key := ""
	return parseTradeData_urlParam(&this.Base,
		orderId_key, channelUserId_key, channelOrderId_key,
		amount_key, discount_key)
}

func (this *Test) CheckChannelRet(params ...interface{}) (err error) {
	tradeState := this.urlParams.Get("result") == "0"
	tradeFailDesc := `urlParam(result)!="0"`
	return this.Base.CheckChannelRet(tradeState, tradeFailDesc)
}

func (this *Test) CheckSign(params ...interface{}) (err error) {
	return
}

func (this *Test) GetResult(params ...interface{}) (ret string) {
	format := `{"code":"%s", "desc:":"%s"}`
	desc := errorDescList[this.lastError]
	succMsg := "0," + desc
	failMsg := "1," + desc
	return this.Base.GetResult(format, succMsg, failMsg)
}
