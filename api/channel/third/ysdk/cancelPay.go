package ysdk

import (
	"net/url"
	"strconv"
	"time"
)

type CancelPayRet struct {
	Ret     int    `json:"ret"`
	Billno  string `json:"billno"`
	Balance int    `json:"balance"`
	ErrCode string `json:"err_code"`
	Msg     string `json:"msg"`
}

type CancelPayParam struct {
	commonParam
	Amt    string `json:"amt"`    //扣游戏币数量，atn不能为0。
	Billno string `json:"billno"` //u9订单号
}

const cancelPayUri = `/mpay/cancel_pay_m`

func GetCancelPayUrl(param *CancelPayParam) (ret string) {
	uri := commonUri + cancelPayUri

	urlParams := url.Values{}
	urlParams.Set("appid", param.AppId)
	urlParams.Set("openid", param.OpenId)
	urlParams.Set("openkey", param.OpenKey)
	urlParams.Set("pay_token", param.PayToken)
	urlParams.Set("pf", param.Pf)
	urlParams.Set("pfkey", param.PfKey)
	urlParams.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	urlParams.Set("zoneid", param.ZoneId)
	urlParams.Set("amt", param.Amt)
	urlParams.Set("billno", param.Billno)

	content, sign := sign(&urlParams, uri, param.PayKey)
	ret = getRootUrl(param.Debug) + uri + `?` + content + "&sig=" + sign
	return
}

func GetCancelPayCookie(loginType string) (ret string, err error) {
	if ret, err = GetCommonCookie(loginType); err != nil {
		return
	}
	ret = ret + ";" + cancelPayUri
	return ret, nil
}
