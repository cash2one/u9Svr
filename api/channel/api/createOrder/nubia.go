package createOrder

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	// "strings"
	"strconv"
	"time"
	"u9/models"
	"u9/tool"
)

type Nubia struct {
	Cr
	channelExt NubiaExt
}

type NubiaExt struct{
	Amount   string   `json:"amount"`
	Number   string   `json:"number"`
	ProductDes  string  `json:"product_des"`
	ProductName  string  `json:"product_name"`
	Uid  string  `json:"uid"`
}

func (this *Nubia) Prepare(lr *models.LoginRequest, orderId, extParamStr string,
	channelParams *map[string]interface{}, ctx *context.Context) (err error) {

	if err = this.Cr.Initial(lr, orderId, nil, nil, extParamStr, channelParams, ctx); err != nil {
		return
	}

	beego.Trace("extParamStr:",extParamStr)

	if err = json.Unmarshal([]byte(extParamStr), &this.channelExt); err != nil{
		beego.Trace(err)
		return
	}

	appid :=  this.parsePayKey("NUBIA_APPID")
	appsecret := this.parsePayKey("NUBIA_APPSECRET_KEY")
	t := strconv.FormatInt(time.Now().Unix(), 10)

	format := "amount=%s&app_id=%s&cp_order_id=%s&data_timestamp=%s&number=%s&product_des=%s&product_name=%s&uid=%s:%s:%s"
	context  := fmt.Sprintf(format,this.channelExt.Amount,appid,orderId,t,"1",this.channelExt.ProductDes,
		this.channelExt.ProductName,this.channelExt.Uid,appid,appsecret)
	sign := tool.Md5([]byte(context))
	beego.Trace("context:",context," #sign:"+sign)
	resultFormat := `{"sign":"%s","data_timestamp":"%s"}`

	this.Result = fmt.Sprintf(resultFormat,sign,t)

	return nil
}

func (this *Nubia) InitParam() (err error) {
	return nil
}

func (this *Nubia) GetResponse() (err error) {
	return nil
}

func (this *Nubia) ParseChannelRet() (err error) {
	return nil
}

func (this *Nubia) GetResult() (ret string) {
	format := "getResult: result:%s"
	msg := fmt.Sprintf(format, this.Result)
	beego.Trace(msg)
	return this.Result
}
