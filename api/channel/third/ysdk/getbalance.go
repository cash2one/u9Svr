package ysdk

import (
	"net/url"
	"strconv"
	"time"
)

type GetBalanceRet struct {
	Ret        int    `json:"ret"`
	Balance    int    `json:"balance"`
	GenBalance int    `json:"gen_balance"`
	FirstSave  int    `json:"first_save"`
	SaveAmt    int    `json:"save_amt"`
	GenExpire  int    `json:"gen_expire"`
	ErrCode    string `json:"err_code"`
	Msg        string `json:"msg"`
}

type GetBalanceParam struct {
	commonParam
}

const getBalanceUri = `/mpay/get_balance_m`

func GetGetBalanceUrl(param *GetBalanceParam) (ret string) {
	uri := commonUri + getBalanceUri

	urlParams := url.Values{}
	urlParams.Set("appid", param.AppId)
	urlParams.Set("openid", param.OpenId)
	urlParams.Set("openkey", param.OpenKey)
	urlParams.Set("pay_token", param.PayToken)
	urlParams.Set("pf", param.Pf)
	urlParams.Set("pfkey", param.PfKey)
	urlParams.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	urlParams.Set("zoneid", param.ZoneId)

	content, sign := sign(&urlParams, uri, param.PayKey)
	ret = getRootUrl(param.Debug) + uri + `?` + content + "&sig=" + sign
	return
}

func GetGetBalanceCookie(loginType string) (ret string, err error) {
	if ret, err = GetCommonCookie(loginType); err != nil {
		return
	}
	ret = ret + ";" + getBalanceUri
	return ret, nil
}
