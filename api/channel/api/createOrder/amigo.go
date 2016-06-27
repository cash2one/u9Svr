package createOrder

import (
	"crypto"
	"crypto/rsa"
	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"strconv"
	"time"
	"u9/api/common"
	"u9/models"
	"u9/tool"
)

type amigoChannelRet struct {
	Status     string `json:"status"`
	Desc       string `json:"description"`
	OrderNo    string `json:"order_no"`
	ApiKey     string `json:"api_key"`
	OutOrderNo string `json:"out_order_no"`
	SubmitTime string `json:"submit_time"`
}

type amigoExtParam struct {
	ApiKey      string `json:"api_key"`
	Subject     string `json:"subject"`
	OutOrderNo  string `json:"out_order_no"`
	DeliverType string `json:"deliver_type"`
	DealPrice   string `json:"deal_price"`
	TotalFee    string `json:"total_fee"`
	SubmitTime  string `json:"submit_time"`
	NotifyUrl   string `json:"notify_url"`
	Sign        string `json:"sign"`
	PlayerId    string `json:"player_id"`
}

type Amigo struct {
	Cr
	privateKey    string
	rsaPrivateKey *rsa.PrivateKey
}

func (this *Amigo) Prepare(lr *models.LoginRequest, orderId, extParamStr string,
	channelParams *map[string]interface{}, ctx *context.Context) (err error) {

	beego.Trace(channelParams)

	this.privateKey = (*channelParams)["AMIGO_PRIVATEKEY"].(string)

	if err = this.initRsaPrivateKey(); err != nil {
		beego.Error(err)
		return err
	}

	if err = this.Cr.Initial(lr, orderId, new(amigoChannelRet), new(amigoExtParam),
		extParamStr, channelParams, ctx); err != nil {
		beego.Error(err)
		return err
	}

	this.Method = "POST"
	this.IsHttps = true
	this.Url = "https://pay.gionee.com/order/create"

	return nil
}

func (this *Amigo) initRsaPrivateKey() (err error) {
	this.rsaPrivateKey, err = tool.ParsePkCS8PrivateKeyWithStr(this.privateKey)
	if err != nil {
		beego.Error(err)
		return err
	}
	return nil
}

func (this *Amigo) InitParam() (err error) {

	if err = this.Cr.InitParam(); err != nil {
		beego.Error(err)
		return
	}

	extParam := this.extParam.(*amigoExtParam)
	extParam.ApiKey = (*this.channelParams)["AMIGO_APIKEY"].(string)
	extParam.OutOrderNo = this.orderId
	extParam.DeliverType = "1"
	extParam.SubmitTime = time.Now().Format(common.CommonTimeLayout)

	extParam.NotifyUrl = "http://" + this.ctx.Input.Host() + "/api/channelPayNotify/" +
		strconv.Itoa(this.lr.ProductId) + "/" + strconv.Itoa(this.lr.ChannelId)
	extParam.PlayerId = this.lr.ChannelUserid

	context := extParam.ApiKey + extParam.DealPrice + extParam.DeliverType + extParam.NotifyUrl +
		extParam.OutOrderNo + extParam.Subject + extParam.SubmitTime + extParam.TotalFee

	extParam.Sign, err = tool.RsaPKCS1V15Sign(this.rsaPrivateKey, crypto.MD5SHA1, context)
	if err != nil {
		beego.Error(err)
		return
	}
	this.Req.JSONBody(extParam)
	return
}

func (this *Amigo) ParseChannelRet() (err error) {
	if err = this.Cr.ParseChannelRet(); err != nil {
		beego.Error(err)
		return
	}

	channelRet := this.channelRet.(*amigoChannelRet)
	if channelRet.Status != "200010000" {
		err = errors.New("status is failure")
		beego.Error(err)
		return
	}
	this.channelOrderId = channelRet.OrderNo
	return
}

func (this *Amigo) GetResult() (ret string) {
	extParam := this.extParam.(*amigoExtParam)
	return extParam.SubmitTime
}
