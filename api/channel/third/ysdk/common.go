package ysdk

import (
	"encoding/base64"
	"errors"
	"github.com/astaxie/beego"
	"net/url"
	"sort"
	"u9/tool"
)

type commonParam struct {
	Debug     bool   `json:"debug"`
	LoginType string `json:"loginType"`
	OpenId    string `json:"open_id"`
	OpenKey   string `json:"open_key"`
	PayToken  string `json:"pay_token"`
	Pf        string `json:"pf"`
	PfKey     string `json:"pf_key"`
	ZoneId    string `json:"zoneid"`
	AppId     string `json:"app_id"`
	PayKey    string `json:"paykey"`
}

const commonUri = `/v3/r`

func GetParamName(loginType string) (appId, appKey, payKey string, err error) {
	if loginType == "QQ" {
		appKey = "QQ_APP_KEY"
	} else if loginType == "WX" {
		appKey = "WX_APP_KEY"
	} else {
		msg := "getParamName: login type must in (QQ, WX)"
		err = errors.New(msg)
		beego.Error(err)
		return
	}
	appId = "QQ_APP_ID"
	payKey = "PAY_KEY"
	return
}

func GetCommonCookie(loginType string) (ret string, err error) {
	ret = ""
	if loginType == "QQ" {
		ret = "session_id=openid;session_type=kp_actoken"
	} else if loginType == "WX" {
		ret = "session_id=hy_gameid;session_type=wc_actoken"
	} else {
		msg := "getCommonCookie: login type must in (QQ, WX)"
		err = errors.New(msg)
		beego.Error(err)
		return
	}
	return ret, nil
}

func getRootUrl(debug bool) (ret string) {
	if debug {
		ret = `https://ysdktest.qq.com`
	} else {
		ret = `https://ysdk.qq.com`
	}
	return
}

func sign(urlParams *url.Values, uri, payKey string) (content, sign string) {
	excludeItems := []string{}
	sorter := tool.NewUrlValuesSorter(urlParams, &excludeItems)
	sort.Sort(sorter)
	content = sorter.DefaultBody()

	signContent := `GET&` + url.QueryEscape(uri) + `&` + url.QueryEscape(content)

	method := "GET"
	signContent = url.QueryEscape(method) + `&` + url.QueryEscape(uri) + `&` + url.QueryEscape(content)

	signResult := tool.HmacSHA1Encrypt(signContent, payKey+`&`)
	sign = base64.StdEncoding.EncodeToString(signResult)
	return
}
