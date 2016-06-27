package createOrder

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"u9/models"
	"u9/tool"
)

type CoolPad struct {
	Cr
}

type coolPadExtParam struct {
	Appid     string  `json:"appid"`
	Waresid   int     `json:"waresid"`
	Cporderid string  `json:"cporderid"`
	Price     float64 `json:"price"`
	Appuserid string  `json:"appuserid"`
	Notifyurl string  `json:"notifyurl"`
}

func (this *CoolPad) Prepare(lr *models.LoginRequest, orderId, extParamStr string,
	channelParams *map[string]interface{}, ctx *context.Context) (err error) {

	if err = this.Cr.Initial(lr, orderId, nil,
		new(coolPadExtParam), extParamStr, channelParams, ctx); err != nil {
		return
	}

	this.parsePayKey("COOLPAD_PRIVATEKEY")

	extParam := this.extParam.(*coolPadExtParam)
	extParam.Cporderid = this.orderId

	var enbyte []byte
	if enbyte, err = json.Marshal(extParam); err != nil {
		format := "prepare: err:%v"
		msg := fmt.Sprintf(format, err)
		beego.Error(msg)
		return
	}

	content := string(enbyte)
	beego.Trace(content)
	beego.Trace(this.payKey)
	if this.Result, err = tool.IapppaySign(content, this.payKey); err != nil {
		format := "prepare: IapppayVerify:%v"
		msg := fmt.Sprintf(format, err)
		beego.Error(msg)
		return
	}

	beego.Trace(this.Result)
	return nil
}

func (this *CoolPad) InitParam() (err error) {
	return nil
}

func (this *CoolPad) GetResponse() (err error) {
	return nil
}

func (this *CoolPad) ParseChannelRet() (err error) {
	return nil
}

func (this *CoolPad) GetResult() (ret string) {
	format := "getResult: result:%s"
	msg := fmt.Sprintf(format, this.Result)
	beego.Trace(msg)
	return this.Result
}
