package createOrder

import (
	"errors"
	"github.com/astaxie/beego/context"
	. "u9/api/channel/third/ysdk"
	"u9/models"
)

type Tencent struct {
	Cr
}

func (this *Tencent) Prepare(lr *models.LoginRequest, orderId, extParamStr string,
	channelParams *map[string]interface{}, ctx *context.Context) (err error) {
	extParam := new(GetBalanceParam)
	if err = this.Cr.Initial(lr, orderId,
		new(GetBalanceRet), extParam, extParamStr, channelParams, ctx); err != nil {
		return
	}

	var key_appId, key_payKey string
	if key_appId, _, key_payKey, err = GetParamName(extParam.LoginType); err != nil {
		return
	}

	extParam.AppId = this.parseAppId(key_appId)
	extParam.PayKey = this.parsePayKey(key_payKey)

	this.Url = GetGetBalanceUrl(extParam)

	this.IsHttps = true
	return nil
}

func (this *Tencent) InitParam() (err error) {
	this.Cr.InitParam()
	extParam := this.extParam.(*GetBalanceParam)
	cookie := ""
	if cookie, err = GetGetBalanceCookie(extParam.LoginType); err != nil {
		return err
	}
	this.Req.Header("cookie", cookie)
	return nil
}

func (this *Tencent) ParseChannelRet() (err error) {
	if err = this.Cr.ParseChannelRet(); err != nil {
		return
	}

	channelRet := this.channelRet.(*GetBalanceRet)
	if channelRet.Ret != 0 {
		err = errors.New("parseChannelRet: channelRet.Ret!=0")
		return
	}
	return
}

func (this *Tencent) GetResult() (ret string) {
	return this.Result
}
