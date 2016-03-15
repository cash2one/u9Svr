package createOrder

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"strconv"
	"time"
	"u9/models"
	"u9/tool"
)

type MeizuUrlParam struct {
	AppId           string `json:"app_id"`
	CpOrderId       string `json:"cp_order_id"`
	UId             string `json:"uid"`
	ProductId       string `json:"product_id"`
	ProductSubject  string `json:"product_subject"`
	ProductBody     string `json:"product_body"`
	ProductUnit     string `json:"product_unit"`
	BuyAmount       string `json:"buy_amount"`
	ProductPerPrice string `json:"product_per_price"`
	TotalPrice      string `json:"total_price"`
	CreateTime      string `json:"create_time"`
	PayType         string `json:"pay_type"`
	UserInfo        string `json:"user_info"`
	Sign            string `json:"sign"`
	SignType        string `json:"sign_type"`
}

type Meizu struct {
	Cr
	channelParams *map[string]interface{}
	urlParam      MeizuUrlParam
	appSecret     string
}

func CoNewMeizu(lr *models.LoginRequest, orderId, host, urlJsonParam string, channelParams *map[string]interface{}) *Meizu {
	ret := new(Meizu)
	ret.Init(lr, orderId, host, urlJsonParam, channelParams)
	return ret
}

func (this *Meizu) Init(lr *models.LoginRequest, orderId, host, urlJsonParam string, channelParams *map[string]interface{}) {
	this.Cr.Init(lr, orderId, host, urlJsonParam)
	this.channelParams = channelParams
}

func (this *Meizu) InitParam() (err error) {
	if err = json.Unmarshal([]byte(this.urlJsonParam), &this.urlParam); err != nil {
		beego.Trace(err, ":", this.urlJsonParam)
		return
	}

	this.appSecret = (*this.channelParams)["MEIZU_APPSECRET"].(string)

	return
}

func (this *Meizu) ParseChannelRet() (err error) {
	return
}

func (this *Meizu) Response() (err error) {
	return
}

func (this *Meizu) GetResult() (ret string) {
	type flymeRet struct {
		CreateTime string `json:"create_time"`
		Sign       string `json:"sign"`
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
