package createOrder

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"strconv"
	"time"
	"u9/models"
	"u9/tool"
)

type meizuExtParam struct {
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

type meizuChannelRet struct {
	CreateTime string `json:"create_time"`
	Sign       string `json:"sign"`
}

type Meizu struct {
	Cr
}

func (this *Meizu) Prepare(lr *models.LoginRequest, orderId, extParamStr string,
	channelParams *map[string]interface{}, ctx *context.Context) (err error) {

	if err = this.Cr.Initial(lr, orderId, new(meizuChannelRet), new(meizuExtParam),
		extParamStr, channelParams, ctx); err != nil {
		beego.Error(err)
		return err
	}

	this.parseAppKey("MEIZU_APPSECRET")
	return nil
}

func (this *Meizu) InitParam() (err error) {
	return nil
}

func (this *Meizu) GetResponse() (err error) {
	return nil
}

func (this *Meizu) ParseChannelRet() (err error) {
	return nil
}

func (this *Meizu) GetResult() (ret string) {
	format := `app_id=%s&buy_amount=%s&cp_order_id=%s&create_time=%s&pay_type=%s&` +
		`product_body=%s&product_id=%s&product_per_price=%s&product_subject=%s&` +
		`product_unit=%s&total_price=%s&uid=%s&user_info=%s:%s`

	extParam := this.extParam.(*meizuExtParam)
	createTime := strconv.FormatInt(time.Now().Unix(), 10)
	context := fmt.Sprintf(format,
		extParam.AppId, extParam.BuyAmount, this.orderId, createTime,
		extParam.PayType, extParam.ProductBody, extParam.ProductId,
		extParam.ProductPerPrice, extParam.ProductSubject, extParam.ProductUnit,
		extParam.TotalPrice, extParam.UId, extParam.UserInfo, this.appKey)

	sign := tool.Md5([]byte(context))

	channelRet := this.channelRet.(*meizuChannelRet)
	channelRet.CreateTime = createTime
	channelRet.Sign = sign

	data, _ := json.Marshal(channelRet)
	ret = string(data)
	return ret
}
