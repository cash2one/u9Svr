package ysdk

import (
	"net/url"
	"strconv"
	"time"
)

type PayRet struct {
	Ret        int    `json:"ret"`
	Billno     string `json:"billno"`
	Balance    int    `json:"balance"`
	GenBalance int    `json:"gen_balance"`
	UsedGenAmt int    `json:"used_gen_amt"`
	ErrCode    string `json:"err_code"`
	Msg        string `json:"msg"`
}

type PayParam struct {
	commonParam
	Amt           string `json:"amt"`    //扣游戏币数量，atn不能为0。
	Billno        string `json:"billno"` //u9订单号
	ChannelUserId string `json:"channelUserId"`
	PayAmount     string `json:"payAmount"`
}

const payUri = `/mpay/pay_m`

func GetPayUrl(param *PayParam) (ret string) {
	uri := commonUri + payUri

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

func GetPayCookie(loginType string) (ret string, err error) {
	if ret, err = GetCommonCookie(loginType); err != nil {
		return
	}
	ret = ret + ";" + payUri
	return ret, nil
}
