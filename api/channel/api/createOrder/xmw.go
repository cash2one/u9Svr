package createOrder

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"strconv"
	"time"
	"u9/models"
	"u9/tool"
)

type xmwChannelRet struct {
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

type xmwExtParam struct {
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
}

func (this *Xmw) Prepare(lr *models.LoginRequest, orderId, extParamStr string,
	channelParams *map[string]interface{}, ctx *context.Context) (err error) {

	if err = this.Cr.Initial(lr, orderId, new(xmwChannelRet), new(xmwExtParam),
		extParamStr, channelParams, ctx); err != nil {
		beego.Error(err)
		return err
	}

	this.Method = "POST"

	tiemStamp := strconv.FormatInt(time.Now().Unix(), 10)
	clientSecret := (*this.channelParams)["XMWAPPSECRET"].(string)

	extParam := this.extParam.(*xmwExtParam)
	format := "amount=%s&app_order_id=%s&app_user_id=%s&notify_url=%s&timestamp=%s&client_secret=%s"
	content := fmt.Sprintf(format, extParam.Amount, this.orderId, extParam.AppUserId,
		extParam.NotifyUrl, tiemStamp, clientSecret)
	//beego.Trace(content)
	sign := tool.Md5([]byte(content))

	format = "http://open.xmwan.com/v2/purchases?" + "amount=%s&app_order_id=%s&app_user_id=%s" +
		"&notify_url=%s&timestamp=%s&sign=%s&access_token=%s&client_id=%s&client_secret=%s"
	clientId := (*this.channelParams)["XMWAPPID"].(string)
	this.Url = fmt.Sprintf(format, extParam.Amount, this.orderId, extParam.AppUserId,
		extParam.NotifyUrl, tiemStamp, sign, extParam.AccessToken, clientId, clientSecret)
	//beego.Trace(this.Url)

	return nil
}

func (this *Xmw) ParseChannelRet() (err error) {
	if err = this.Cr.ParseChannelRet(); err != nil {
		beego.Error(err)
		return
	}

	channelRet := this.channelRet.(*xmwChannelRet)
	if channelRet.Error != "" {
		err = errors.New(channelRet.Error)
		beego.Error(err)
		return
	}
	this.channelOrderId = channelRet.Serial
	return
}

func (this *Xmw) GetResult() (ret string) {
	return this.Result
}
