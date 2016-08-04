package createOrder

import (
	// "encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"strings"
	"u9/models"
	"u9/tool"
	"encoding/json"
)

type CoolPad struct {
	Cr
	coolPadJson CoolPadJson
}

type CoolPadJson struct {
	ProductId 	 string 	`json:"productId"`
	DataJson     string     `json:"channelData"`
}

func (this *CoolPad) Prepare(lr *models.LoginRequest, orderId, extParamStr string,
	channelParams *map[string]interface{}, ctx *context.Context) (err error) {

	beego.Trace("extParamStr:",extParamStr)
	if err = json.Unmarshal([]byte(extParamStr), &this.coolPadJson); err != nil {
		//旧版本传入字符串，会导致json 失败
		if err = this.Cr.Initial(lr, orderId, nil, nil, extParamStr, channelParams, ctx); err != nil {
			return
		}
		this.parsePayKey("COOLPAD_PRIVATEKEY")

		content := fmt.Sprintf(extParamStr, this.orderId)
		content = strings.Replace(content, `\`, ``, -1) //去json中的`\`

		if this.Result, err = tool.IapppaySign(content, this.payKey); err != nil {
			format := "prepare: IapppayVerify:%v"
			msg := fmt.Sprintf(format, err)
			beego.Error(msg)
			return
		}
	}else{
		beego.Trace("productId:", this.coolPadJson.ProductId)
		beego.Trace("channelData:",this.coolPadJson.DataJson)

		if err = this.Cr.Initial(lr, orderId, nil, nil, extParamStr, channelParams, ctx); err != nil {
			return
		}
		payId := this.parsePayKey(this.coolPadJson.ProductId)
		this.parsePayKey("COOLPAD_PRIVATEKEY")

		content := strings.Replace(this.coolPadJson.DataJson, `"u9PayId"` ,payId , -1)
		content = strings.Replace(content, `u9OrderId`, this.orderId, -1) //去json中的`\`
		content = strings.Replace(content, `\`, ``, -1) //去json中的`\`
		beego.Trace("content:",content)
		beego.Trace("payKey:",this.payKey)
		var sign  string
		if sign , err = tool.IapppaySign(content, this.payKey); err != nil {
			format := "prepare: IapppayVerify:%v"
			msg := fmt.Sprintf(format, err)
			beego.Error(msg)
			return
		}
		// beego.Trace("sign:",sign)
		content = `{"sign":"%s","payId":"%s"}`
		this.Result = fmt.Sprintf(content,sign,payId)
		beego.Trace("Result:",this.Result)
	}
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
