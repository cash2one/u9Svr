package createOrder

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"net/url"
	"strings"
	"u9/models"
	"u9/tool"
)

type huaweiChannelRet struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
}

type huaweiExtParam struct {
	Content string `json:"content"`
}

type Huawei struct {
	Cr
}

var huaweiRsaPrivateKey *rsa.PrivateKey

func (this *Huawei) prepareRsaPrivateKey() (err error) {
	if huaweiRsaPrivateKey == nil {
		rsaPrivateKeyStr := ""
		ok := false
		if rsaPrivateKeyStr, ok = (*this.channelParams)["HUAWEI_PAY_PRIVATE_KEY"].(string); !ok {
			err = errors.New("getPackageParam is error")
			return
		}
		if huaweiRsaPrivateKey, err = tool.ParsePkCS8PrivateKeyWithStr(rsaPrivateKeyStr); err != nil {
			return
		}
	}
	return nil
}

func (this *Huawei) Prepare(lr *models.LoginRequest, orderId, extParamStr string,
	channelParams *map[string]interface{}, ctx *context.Context) (err error) {

	if err = this.Cr.Initial(lr, orderId, new(huaweiChannelRet), new(huaweiExtParam),
		extParamStr, channelParams, ctx); err != nil {
		beego.Error(err)
		return err
	}

	if err = this.prepareRsaPrivateKey(); err != nil {
		beego.Error(err)
		return err
	}
	return nil
}

func (this *Huawei) InitParam() (err error) {
	return
}

func (this *Huawei) GetResponse() (err error) {
	return
}

func (this *Huawei) ParseChannelRet() (err error) {
	return
}

func (this *Huawei) GetResult() (ret string) {
	var err error
	defer func() {
		if err != nil {
			beego.Error(err)
		}
	}()

	extParam := this.extParam.(*huaweiExtParam)
	channelRet := this.channelRet.(*huaweiChannelRet)

	content := ""
	if content, err = url.QueryUnescape(extParam.Content); err != nil {
		channelRet.Status = 1
		channelRet.Msg = "content urlDecode error."
		return
	}
	content = strings.Replace(content, "orderId", this.orderId, -1)
	beego.Trace(content)
	if channelRet.Msg, err = tool.RsaPKCS1V15Sign(huaweiRsaPrivateKey, content); err != nil {
		channelRet.Status = 2
		channelRet.Msg = "rsa sign error."
		return
	}
	channelRet.Status = 0
	data, _ := json.Marshal(channelRet)
	ret = string(data)
	beego.Trace(ret)
	return
}
