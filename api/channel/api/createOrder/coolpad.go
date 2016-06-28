package createOrder

import (
	//"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"strings"
	"u9/models"
	"u9/tool"
)

type CoolPad struct {
	Cr
}

func (this *CoolPad) Prepare(lr *models.LoginRequest, orderId, extParamStr string,
	channelParams *map[string]interface{}, ctx *context.Context) (err error) {

	if err = this.Cr.Initial(lr, orderId, nil, nil, extParamStr, channelParams, ctx); err != nil {
		return
	}

	this.parsePayKey("COOLPAD_PRIVATEKEY")

	content := fmt.Sprintf(this.extParamStr, this.orderId)
	content = strings.Replace(content, `\`, ``, -1) //去json中的`\`

	if this.Result, err = tool.IapppaySign(content, this.payKey); err != nil {
		format := "prepare: IapppayVerify:%v"
		msg := fmt.Sprintf(format, err)
		beego.Error(msg)
		return
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
