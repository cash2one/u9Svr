package createOrder

import (
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
}

const amigoRsaPrivateKeyStr = `MIICdwIBADANBgkqhkiG9w0BAQEFAASCAmEwggJdAgEAAoGBAJs5MKkLs2c+7xORPLvx+S6roUqThA4pcF03bl2+BUjpz7WJMaFOqpkTedTYcHANZN/ztHb/L5VddjB0XaTsN1f6fqukPdrIeuVKN3ioDJlHrStCgf55pckXtMZiAqNwvsiZMVZNHf0QzVhPNX9dgEGO34B2lmzhyBLvl2cJEKPRAgMBAAECgYArZTe3avADA1Mvc0E5ghKZ+52iDc+zbd5eESsnxgIQOl25cNuRLz7+gLVkzgXRLc1v1uAzDHCvH2v1a/LqMqfd3mZtzShrLnccgCCY64XC3xxSSBsVReYWV5uJP9g7PU2P5eatS1aUopWkNYSq25atE7dpxlMJ4T4tcnDQvav6gQJBAOKq8DD0Cx9mBmS95PazMaHzaM56Wb/8tEHoZDsTaJ5nW+RW7SnYPmHXlKVs5OpQTF3FU7KDgj/QjGkvW9Oyh1kCQQCvT3Ce4PVfjiobFczdmpgQJlUn+UK86TX0JdmRWsXvpOXo0m4HIk/BV/pyPtwBGKSJ0b/s8gXf603YR8BWlWk5AkEAiieYOK42rU+ZLAQWL0uvP7/FrLwkQgF7uQQ1O1CsHohvGPDmou+brjUg8+c4a5y/vxPL3O2NEOpC+sWT2adiGQJABpB88RYPWhKitPzt/OZLB1/IFIUa4KQC5y97pBu4Ca8tBLjMcevw/JZkxF5iMpBPqPF3tFGjsqzG73BQXW2e0QJBAI5CpS0S7g6IgvhGSqavonFR3Pkgkbdn5qKIATZCywkBeMx0QHTO3EJq/yLWacsi4Na5l6SNcX3Tde6hwRpVuYw=`

var (
	amigoRsaPrivateKey *rsa.PrivateKey
)

func (this *Amigo) Prepare(lr *models.LoginRequest, orderId, extParamStr string,
	channelParams *map[string]interface{}, ctx *context.Context) (err error) {

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
	if amigoRsaPrivateKey == nil {
		amigoRsaPrivateKey, err = tool.ParsePkCS8PrivateKeyWithStr(amigoRsaPrivateKeyStr)
		if err != nil {
			beego.Error(err)
			return err
		}
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

	extParam.Sign, err = tool.RsaPKCS1V15Sign(amigoRsaPrivateKey, context)
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
