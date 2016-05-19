package createOrder

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"net/url"
	"strconv"
	"time"
	"u9/models"
	"u9/tool"
)

const TencentPayQueryOrgLoc = `/v3/r/mpay/get_balance_m`

type TencentChannelRet struct {
	Ret        int    `json:"ret"`
	Balance    int    `json:"balance"`
	GenBalance int    `json:"gen_balance"`
	FirstSave  int    `json:"first_save"`
	SaveAmt    int    `json:"save_amt"`
	GenExpire  int    `json:"gen_expire"`
	ErrCode    string `json:"err_code"`
	Msg        string `json:"msg"`
}

type TencentExtParam struct {
	Debug     bool   `json:"debug"`
	LoginType string `json:"loginType"`
	PayToken  string `json:"pay_token"`
	OpenKey   string `json:"open_key"`
	OpenId    string `json:"open_id"`
	Pf        string `json:"pf"`
	PfKey     string `json:"pf_key"`
	ZoneId    string `json:"zoneid"`
}

type Tencent struct {
	Cr
}

func GetTencentPayQueryUrl(debug bool,
	openId, openKey, payToken, pf, pfKey, zoneId, appId, payKey string) (ret string) {
	if debug {
		ret = `https://ysdktest.qq.com` + TencentPayQueryOrgLoc + `?`
	} else {
		ret = `https://ysdk.qq.com` + TencentPayQueryOrgLoc + `?`
	}

	method := `GET`
	format := `appid=%s&openid=%s&openkey=%s&pay_token=%s&pf=%s&pfkey=%s&ts=%s&zoneid=%s`
	timeStamp := strconv.FormatInt(time.Now().Unix(), 10)
	content := fmt.Sprintf(format, appId, openId, openKey, payToken, pf, pfKey, timeStamp, zoneId)

	signContent := url.QueryEscape(method) + `&` +
		url.QueryEscape(TencentPayQueryOrgLoc) + `&` + url.QueryEscape(content)

	signResult := tool.HmacSHA1Encrypt(signContent, payKey+`&`)
	encodeSign := base64.StdEncoding.EncodeToString(signResult)

	content = content + "&sig=" + encodeSign
	ret = ret + content

	beego.Trace(ret)
	return
}

func GetTencentPayQueryCookie(loginType string) (ret string, err error) {
	ret = ""
	if loginType == "QQ" {
		ret = "session_id=openid;session_type=kp_actoken"
	} else if loginType == "WX" {
		ret = "session_id=hy_gameid;session_type=wc_actoken"
	} else {
		err = errors.New("login type is error, must in (QQ, WX)")
		beego.Error(err)
		return
	}
	ret = ret + ";" + TencentPayQueryOrgLoc
	beego.Trace(ret)
	return ret, nil
}

func GetTencentPayParamName(loginType string) (appId, appKey, payKey string, err error) {
	if loginType == "QQ" {
		appKey = "QQ_APP_KEY"
	} else if loginType == "WX" {
		appKey = "WX_APP_KEY"
	} else {
		err = errors.New("login type is error, must in (QQ, WX)")
		beego.Error(err)
		return
	}
	appId = "QQ_APP_ID"
	payKey = "PAY_KEY"
	return
}

func (this *Tencent) Prepare(lr *models.LoginRequest, orderId, extParamStr string,
	channelParams *map[string]interface{}, ctx *context.Context) (err error) {

	if err = this.Cr.Initial(lr, orderId, new(TencentChannelRet), new(TencentExtParam),
		extParamStr, channelParams, ctx); err != nil {
		beego.Error(err)
		return err
	}

	this.IsHttps = true
	extParam := this.extParam.(*TencentExtParam)

	var appId, appKey, payKey string
	if appId, appKey, payKey, err = GetTencentPayParamName(extParam.LoginType); err != nil {
		beego.Error("extParamStr:" + extParamStr)
		return
	}
	this.parseAppId(appId)
	this.parseAppKey(appKey)
	this.parsePayKey(payKey)

	this.Url = GetTencentPayQueryUrl(extParam.Debug,
		extParam.OpenId, extParam.OpenKey, extParam.PayToken, extParam.Pf,
		extParam.PfKey, extParam.ZoneId, this.appId, this.payKey)

	return nil
}

func (this *Tencent) InitParam() (err error) {
	this.Cr.InitParam()
	extParam := this.extParam.(*TencentExtParam)
	cookie := ""
	if cookie, err = GetTencentPayQueryCookie(extParam.LoginType); err != nil {
		return err
	}
	this.Req.Header("cookie", cookie)
	return nil
}

func (this *Tencent) ParseChannelRet() (err error) {
	if err = this.Cr.ParseChannelRet(); err != nil {
		beego.Error(err)
		return
	}

	channelRet := this.channelRet.(*TencentChannelRet)
	if channelRet.Ret != 0 {
		err = errors.New("status is failure")
		beego.Error(err)
		return
	}

	return
}

func (this *Tencent) GetResult() (ret string) {
	return this.Result
}
