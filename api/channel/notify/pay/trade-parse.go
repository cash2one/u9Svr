package channelPayNotify

import (
	"fmt"
	"github.com/astaxie/beego"
	"reflect"
)

func parseTradeData_urlParam(pn *Base, paramKeys ...interface{}) (err error) {
	parseMethod := "parseTradeData_urlParam"
	orderId_key := paramKeys[0].(string)

	if orderId_key != "" {
		pn.orderId = pn.urlParams.Get(orderId_key)
		if !checkOrderId(pn, parseMethod, paramKeys...) {
			return
		}
	}

	channelUserId_key := paramKeys[1].(string)
	if channelUserId_key != "" {
		pn.channelUserId = pn.urlParams.Get(channelUserId_key)
		if !checkChannelUserId(pn, parseMethod, paramKeys...) {
			return
		}
	}

	channelOrderId_key := paramKeys[2].(string)
	if channelOrderId_key != "" {
		pn.channelOrderId = pn.urlParams.Get(channelOrderId_key)
		if !checkChannelOrderId(pn, parseMethod, paramKeys...) {
			return
		}
	}

	//payAmount
	amount := ""
	amount_key := paramKeys[3].(string)
	if amount_key != "" {
		amount = pn.urlParams.Get(amount_key)
	}

	discount := ""
	discount_key := paramKeys[4].(string)
	if discount_key != "" {
		discount = pn.urlParams.Get(discount_key)
	}

	if err = checkPayAmount(pn, parseMethod, amount, discount, paramKeys...); err != nil {
		return
	}

	return
}

func parseTradeData_stringMap(pn *Base, paramKeys ...interface{}) (err error) {
	parseMethod := "parseTradeData_stringMap"
	var retInterface interface{}
	channelTradeData := pn.channelTradeData.(map[string]interface{})

	orderId_key := paramKeys[0].(string)
	if orderId_key != "" {
		if retInterface, err = parseJsonMapParams(&channelTradeData, orderId_key); err != nil {
			beego.Error(err)
		} else {
			pn.orderId = retInterface.(string)
		}

		if !checkOrderId(pn, parseMethod, paramKeys...) {
			return
		}
	}

	channelUserId_key := paramKeys[1].(string)
	if channelUserId_key != "" {
		if retInterface, err = parseJsonMapParams(&channelTradeData, channelUserId_key); err != nil {
			beego.Error(err)
		} else {
			pn.channelUserId = retInterface.(string)
		}

		if !checkChannelUserId(pn, parseMethod, paramKeys...) {
			return
		}
	}

	channelOrderId_key := paramKeys[2].(string)
	if channelOrderId_key != "" {
		if retInterface, err = parseJsonMapParams(&channelTradeData, channelOrderId_key); err != nil {
			beego.Error(err)
		} else {
			pn.channelOrderId = retInterface.(string)
		}

		if !checkChannelOrderId(pn, parseMethod, paramKeys...) {
			return
		}
	}

	//payAmount
	amount := ""
	amount_key := paramKeys[3].(string)
	if amount_key != "" {
		if retInterface, err = parseJsonMapParams(&channelTradeData, amount_key); err != nil {
			beego.Error(err)
		} else {
			amount = retInterface.(string)
		}
	}

	discount := ""
	discount_key := paramKeys[4].(string)
	if discount_key != "" {
		if retInterface, err = parseJsonMapParams(&channelTradeData, discount_key); err != nil {
			beego.Error(err)
		} else {
			discount = retInterface.(string)
		}
	}

	if err = checkPayAmount(pn, parseMethod, amount, discount, paramKeys...); err != nil {
		return
	}

	return
}

func parseTradeData_customStruct(pn *Base, paramKeys ...interface{}) (err error) {
	parseMethod := "parseTradeData_customStruct"
	var retInterface reflect.Value
	channelTradeData := pn.channelTradeData

	orderId_key := paramKeys[0].(string)
	if orderId_key != "" {
		if retInterface, err = parseCustomStructParams(interface{}(channelTradeData), orderId_key); err != nil {
			beego.Error(err)
		}
		pn.orderId = fmt.Sprintf("%v", retInterface)
		if !checkOrderId(pn, parseMethod, paramKeys...) {
			return
		}
	}

	channelUserId_key := paramKeys[1].(string)
	if channelUserId_key != "" {
		if retInterface, err = parseCustomStructParams(interface{}(channelTradeData), channelUserId_key); err != nil {
			beego.Error(err)
		}
		pn.channelUserId = fmt.Sprintf("%v", retInterface)
		if !checkChannelUserId(pn, parseMethod, paramKeys...) {
			return
		}
	}

	channelOrderId_key := paramKeys[2].(string)
	if channelOrderId_key != "" {
		if retInterface, err = parseCustomStructParams(interface{}(channelTradeData), channelOrderId_key); err != nil {
			beego.Error(err)
		}

		pn.channelOrderId = fmt.Sprintf("%v", retInterface)
		if !checkChannelOrderId(pn, parseMethod, paramKeys...) {
			return
		}
	}

	//payAmount
	amount := ""
	amount_key := paramKeys[3].(string)
	if amount_key != "" {
		if retInterface, err = parseCustomStructParams(interface{}(channelTradeData), amount_key); err != nil {
			beego.Error(err)
		}
		amount = fmt.Sprintf("%v", retInterface)
	}

	discount := ""
	discount_key := paramKeys[4].(string)
	if discount_key != "" {
		if retInterface, err = parseCustomStructParams(interface{}(channelTradeData), discount_key); err != nil {
			beego.Error(err)
		}
		discount = fmt.Sprintf("%v", retInterface)
	}

	if err = checkPayAmount(pn, parseMethod, amount, discount, paramKeys...); err != nil {
		return
	}

	return
}
