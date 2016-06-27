package channelPayNotify

import (
	"fmt"
	"github.com/astaxie/beego"
)

func dumpUrlTradeParam(paramKeys ...interface{}) (ret string) {
	format := `orderId_key:%s, channelUserId_key:%s, channelOrderId_key:%s amount_key:%s, discount_key:%s`

	orderId_key := paramKeys[0].(string)
	channelUserId_key := paramKeys[1].(string)
	channelOrderId_key := paramKeys[2].(string)
	amount_key := paramKeys[3].(string)
	discount_key := paramKeys[4].(string)

	ret = fmt.Sprintf(format, orderId_key, channelUserId_key, channelOrderId_key, amount_key, discount_key)
	return
}

func checkOrderId(pn *Base, parseMethod string, paramKeys ...interface{}) (ret bool) {
	if pn.orderId == "" {
		pn.lastError = err_parseOrderId
		beego.Error(parseMethod + ":" + dumpUrlTradeParam(paramKeys...))
		ret = false
	} else {
		ret = true
	}
	return
}

func checkChannelUserId(pn *Base, parseMethod string, paramKeys ...interface{}) (ret bool) {
	if pn.channelUserId == "" {
		pn.lastError = err_parseChannelUserId
		beego.Error(parseMethod + ":" + dumpUrlTradeParam(paramKeys...))
		ret = false
	} else {
		ret = true
	}
	return
}

func checkChannelOrderId(pn *Base, parseMethod string, paramKeys ...interface{}) (ret bool) {
	if pn.channelOrderId == "" {
		pn.lastError = err_parseChannelOrderId
		beego.Error(parseMethod + ":" + dumpUrlTradeParam(paramKeys...))
		ret = false
	} else {
		ret = true
	}
	return
}

func checkPayAmount(pn *Base, parseMethod, amount, discount string, paramKeys ...interface{}) (err error) {
	if err = pn.parsePayAmount(amount, discount); err != nil {
		beego.Error(parseMethod + ":" + dumpUrlTradeParam(paramKeys...))
	}
	return
}
