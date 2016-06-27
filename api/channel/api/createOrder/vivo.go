package createOrder

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"time"
	"u9/models"
	"u9/tool"
)

type vivoChannelRet struct {
	RespCode    string `json:"respCode"`
	RespMsg     string `json:"respMsg"`
	SignMethod  string `json:"signMethod"`
	Signature   string `json:"signature"`
	AccessKey   string `json:"accessKey"`
	OrderNumber string `json:"orderNumber"`
	OrderAmount string `json:"orderAmount"`
}

type vivoExtParam struct {
	Version       string `json:"version"`
	SignMethod    string `json:"signMethod"`
	Signature     string `json:"signature"`
	CpId          string `json:"cpId"`
	AppId         string `json:"appId"`
	CpOrderNumber string `json:"cpOrderNumber"`
	NotifyUrl     string `json:"notifyUrl"`
	OrderTime     string `json:"orderTime"`
	OrderAmount   string `json:"orderAmount"`
	OrderTitle    string `json:"orderTitle"`
	OrderDesc     string `json:"orderDesc"`
	ExtInfo       string `json:"extInfo"`
}

type Vivo struct {
	Cr
}

func (this *Vivo) Prepare(lr *models.LoginRequest, orderId, extParamStr string,
	channelParams *map[string]interface{}, ctx *context.Context) (err error) {

	if err = this.Cr.Initial(lr, orderId, new(vivoChannelRet), new(vivoExtParam),
		extParamStr, channelParams, ctx); err != nil {
		beego.Error(err)
		return err
	}
	this.Method = "POST"
	this.IsHttps = true
	this.Url = "https://pay.vivo.com.cn/vcoin/trade"
	beego.Trace(this.Url)
	return nil
}

func (this *Vivo) InitParam() (err error) {
	if err = this.Cr.InitParam(); err != nil {
		beego.Error(err)
		return
	}
	extParam := this.extParam.(*vivoExtParam)
	extParam.SignMethod = "MD5"
	extParam.CpId = (*this.channelParams)["VIVO_CP_ID"].(string)
	extParam.AppId = (*this.channelParams)["APP_MONITOR_APPID"].(string)
	extParam.OrderTime = time.Now().Format("20060102150405")
	extParam.CpOrderNumber = this.orderId

	cpKey := (*this.channelParams)["VIVO_CP_KEY"].(string)
	format := `appId=%s&cpId=%s&cpOrderNumber=%s&extInfo=%s&notifyUrl=%s&orderAmount=%s&` +
		`orderDesc=%s&orderTime=%s&orderTitle=%s&version=%s&%s`
	content := fmt.Sprintf(format, extParam.AppId, extParam.CpId, extParam.CpOrderNumber,
		extParam.ExtInfo, extParam.NotifyUrl, extParam.OrderAmount, extParam.OrderDesc,
		extParam.OrderTime, extParam.OrderTitle, extParam.Version, tool.Md5([]byte(cpKey)))
	beego.Trace("content:" + content)
	extParam.Signature = tool.Md5([]byte(content))

	this.Req.Param("version", extParam.Version)
	this.Req.Param("signMethod", extParam.SignMethod)
	this.Req.Param("cpId", extParam.CpId)
	this.Req.Param("appId", extParam.AppId)
	this.Req.Param("cpOrderNumber", extParam.CpOrderNumber)
	this.Req.Param("notifyUrl", extParam.NotifyUrl)
	this.Req.Param("orderTime", extParam.OrderTime)
	this.Req.Param("orderAmount", extParam.OrderAmount)
	this.Req.Param("orderTitle", extParam.OrderTitle)
	this.Req.Param("orderDesc", extParam.OrderDesc)
	this.Req.Param("extInfo", extParam.ExtInfo)
	this.Req.Param("signature", extParam.Signature)

	beego.Trace(fmt.Sprintf("%s%+v", "extParam:", extParam))

	return
}

func (this *Vivo) ParseChannelRet() (err error) {
	if err = this.Cr.ParseChannelRet(); err != nil {
		beego.Error(err)
		return
	}

	channelRet := this.channelRet.(*vivoChannelRet)
	if channelRet.RespCode != "200" {
		msg := fmt.Sprintf("%s%+v", "channelRet:", channelRet)
		err = errors.New(msg)
		beego.Error(err)
		return
	}
	this.channelOrderId = channelRet.OrderNumber
	return
}

func (this *Vivo) GetResult() (ret string) {
	return this.Result
}
