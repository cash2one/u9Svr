package createOrder

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"github.com/astaxie/beego"
	"strconv"
	"time"
	"u9/api/common"
	"u9/models"
	"u9/tool"
)

type XmwCoChannelRet struct {
	Serial         string `json:"serial"`
	NotifyUrl      string `json:"notify_url"`
	Amount         string `json:"amount"`
	Cost           string `json:"cost"`
	AppOrderId     string `json:"app_order_id"`
	AppUserId      string `json:"app_user_id"`
	AppSubject     string `json:"app_subject"`
	AppDescription string `json:"app_description"`
	AppExt1        string `json:"app_ext1"`
	AppExt2        string `json:"app_ext2"`
}

type XmwUrlParam struct {
	AccessToken    string `json:"access_token"`
	ClientId       string `json:"client_id"`
	ClientSecret   string `json:"client_secret"`
	AppOrderId     string `json:"app_order_id"`
	AppUserId      string `json:"app_user_id"`
	NotifyUrl      string `json:"notify_url"`
	Amount         int    `json:"amount"`
	Timestamp      int    `json:"timestamp"`
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
	this.Url = "http://open.xmwan.com/v2/purchases"
	this.channelParams = channelParams
}

func (this *Xmw) InitParam() (err error) {
	this.Cr.InitParam()

	if err = json.Unmarshal([]byte(this.urlJsonParam), &this.urlParam); err != nil {
		beego.Trace(err, ":", this.urlJsonParam)
	}

	this.Req.Param("amount", this.urlParam.Amount)
	this.Req.Param("app_order_id", this.orderId)
	this.Req.Param("app_user_id", this.urlParam.AppUserId)
	this.Req.Param("access_token", this.urlParam.AccessToken)
	this.Req.Param("client_id", (*this.channelParams)["XMWAPPID"].(string))
	this.Req.Param("client_secret", (*this.channelParams)["XMWAPPSECRET"].(string))

	this.Req.Param("notify_url", this.urlParam.NotifyUrl)

	this.Req.Param("timestamp", strconv.FormatInt(time.Now().Unix(), 10))
	this.Req.Param("app_subject", this.urlParam.AppSubject)
	this.Req.Param("app_description", this.urlParam.AppDescription)
	this.Req.Param("app_ext1", this.urlParam.AppExt1)
	this.Req.Param("app_ext2", this.urlParam.AppExt2)

	format := "amount=%s&app_order_id=%s&app_user_id=%s&notify_url=%s&timestamp=%s&client_secret=%s"
	this.Req.Param("sign", "&client_secret")
	return
}

func (this *Xmw) ParseChannelRet() (err error) {
	beego.Trace(this.Result)
	if err = json.Unmarshal([]byte(this.Result), &this.channelRet); err != nil {
		beego.Error(err)
		return
	}
	if this.channelRet.Status != "200010000" {
		err = errors.New("status is failure")
		beego.Error(err)
		return
	}
	this.ChannelOrderId = this.channelRet.OrderNo
	return
}

func (this *Xmw) GetResult() (ret string) {

	channelRet := XmwCoChannelRet{
		Serial:         "",
		NotifyUrl:      this.urlParam.NotifyUrl,
		Amount:         strconv.Itoa(this.urlParam.Amount),
		Cost:           strconv.Itoa(this.urlParam.Amount),
		AppOrderId:     this.urlParam.AppOrderId,
		AppUserId:      this.urlParam.AppUserId,
		AppSubject:     this.urlParam.AppSubject,
		AppDescription: this.urlParam.AppDescription,
		AppExt1:        this.urlParam.AppExt1,
		AppExt2:        this.urlParam.AppExt2,
	}
	format := `app_id=%s&buy_amount=%s&cp_order_id=%s&create_time=%s&pay_type=%s&product_body=%s&product_id=%s&product_per_price=%s&product_subject=%s&product_unit=%s&total_price=%s&uid=%s&user_info=%s:%s`

	createTime := strconv.FormatInt(time.Now().Unix(), 10)
	context := fmt.Sprintf(format,
		this.urlParam.AppId, this.urlParam.BuyAmount, this.orderId, createTime,
		this.urlParam.PayType, this.urlParam.ProductBody, this.urlParam.ProductId,
		this.urlParam.ProductPerPrice, this.urlParam.ProductSubject, this.urlParam.ProductUnit,
		this.urlParam.TotalPrice, this.urlParam.UId, this.urlParam.UserInfo, this.appSecret)

	sign := tool.Md5([]byte(context))

	jsonRet := flymeRet{
		CreateTime: createTime,
		Sign:       sign,
	}

	data, _ := json.Marshal(jsonRet)
	ret = string(data)
	return ret
}
