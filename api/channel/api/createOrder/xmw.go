package createOrder

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"strconv"
	"time"
	"u9/models"
	"u9/tool"
)

type XmwCoChannelRet struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
	Serial           string `json:"serial"`
	NotifyUrl        string `json:"notify_url"`
	Amount           int    `json:"amount"`
	Cost             int    `json:"cost"`
	AppOrderId       string `json:"app_order_id"`
	AppUserId        string `json:"app_user_id"`
	AppSubject       string `json:"app_subject"`
	AppDescription   string `json:"app_description"`
	AppExt1          string `json:"app_ext1"`
	AppExt2          string `json:"app_ext2"`
}

type XmwUrlParam struct {
	AccessToken    string `json:"access_token"`
	ClientId       string `json:"client_id"`
	ClientSecret   string `json:"client_secret"`
	AppOrderId     string `json:"app_order_id"`
	AppUserId      string `json:"app_user_id"`
	NotifyUrl      string `json:"notify_url"`
	Amount         string `json:"amount"`
	Timestamp      string `json:"timestamp"`
	Sign           string `json:"sign"`
	AppSubject     string `json:"app_subject"`
	AppDescription string `json:"app_description"`
	AppExt1        string `json:"app_ext1"`
	AppExt2        string `json:"app_ext2"`
}

type Xmw struct {
	Cr
	channelParams *map[string]interface{}
	urlParam      XmwUrlParam
	channelRet    XmwCoChannelRet
}

func CoNewXmw(lr *models.LoginRequest, orderId, host, urlJsonParam string, channelParams *map[string]interface{}) *Xmw {
	ret := new(Xmw)
	ret.Init(lr, orderId, host, urlJsonParam, channelParams)
	return ret
}

func (this *Xmw) Init(lr *models.LoginRequest, orderId, host, urlJsonParam string, channelParams *map[string]interface{}) {
	this.Cr.Init(lr, orderId, host, urlJsonParam)
	this.Method = "POST"
	this.channelParams = channelParams
}

func (this *Xmw) InitParam() (err error) {
	if err = json.Unmarshal([]byte(this.urlJsonParam), &this.urlParam); err != nil {
		beego.Trace(err, ":", this.urlJsonParam)
	}

	tiemStamp := strconv.FormatInt(time.Now().Unix(), 10)
	clientSecret := (*this.channelParams)["XMWAPPSECRET"].(string)
	clientId := (*this.channelParams)["XMWAPPID"].(string)

	format := "amount=%s&app_order_id=%s&app_user_id=%s&notify_url=%s&timestamp=%s&client_secret=%s"
	content := fmt.Sprintf(format,
		this.urlParam.Amount, this.orderId,
		this.urlParam.AppUserId, this.urlParam.NotifyUrl, tiemStamp, clientSecret)
	beego.Trace(content)
	sign := tool.Md5([]byte(content))

	format = "http://open.xmwan.com/v2/purchases?" +
		"amount=%s&app_order_id=%s&app_user_id=%s" +
		"&notify_url=%s&timestamp=%s&sign=%s&access_token=%s&client_id=%s&client_secret=%s"
	this.Url = fmt.Sprintf(format,
		this.urlParam.Amount, this.orderId,
		this.urlParam.AppUserId, this.urlParam.NotifyUrl,
		tiemStamp, sign, this.urlParam.AccessToken, clientId, clientSecret)
	beego.Trace(this.Url)

	this.Cr.InitParam()
	return
}

func (this *Xmw) ParseChannelRet() (err error) {
	beego.Trace(this.Result)
	if err = json.Unmarshal([]byte(this.Result), &this.channelRet); err != nil {
		beego.Error(err)
		return
	}
	if this.channelRet.Error != "" {
		err = errors.New(this.channelRet.Error)
		beego.Error(err)
		return
	}
	this.ChannelOrderId = this.channelRet.Serial
	return
}

func (this *Xmw) GetResult() (ret string) {
	return this.Result
}
