package channelPayNotify

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

//example: ret, err := parseCustomStructParams(interface{}(&data), "TradeData.Desc")
func parseCustomStructParams(data interface{}, param string) (ret reflect.Value, err error) {
	defer func() {
		if errInterface := recover(); errInterface != nil {
			format := `panic:%v, data:%+v, param:%s`
			msg := fmt.Sprintf(format, errInterface, data, param)
			err = errors.New(msg)
		}
	}()

	spiltParams := strings.Split(param, ".")
	ret = reflect.ValueOf(data).Elem()

	for _, paramKey := range spiltParams {
		ret = ret.FieldByName(paramKey)
	}
	return
}

//example: ret, err := parseJsonMapParams(&data,"TradeData.Desc")
func parseJsonMapParams(data *map[string]interface{}, param string) (ret interface{}, err error) {
	curParamValue := *data
	defer func() {
		if errInterface := recover(); errInterface != nil {
			format := `panic:%v, curParamValue:%+v, data:%+v, param:%s`
			msg := fmt.Sprintf(format, errInterface, curParamValue, data, param)
			err = errors.New(msg)
		}
	}()

	spiltParams := strings.Split(param, ".")
	maxIndex := len(spiltParams) - 1

	for index, paramKey := range spiltParams {
		if maxIndex == index {
			ret = curParamValue[paramKey]
		} else {
			curParamValue = curParamValue[paramKey].(map[string]interface{})
		}
	}
	return
}
