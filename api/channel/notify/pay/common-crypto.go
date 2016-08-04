package channelPayNotify

import (
	"crypto"
	"crypto/rsa"
	"fmt"
	"github.com/astaxie/beego"
	"strings"
	"u9/tool"
)

type Crypto struct {
	Base
	signMethod  string
	inputSign   string
	signContent string
}

type MD5 struct {
	Crypto
	signHandleMethod string
}

func (this *MD5) CheckSign(params ...interface{}) (err error) {
	this.signMethod = "md5"

	sign := tool.Md5([]byte(this.signContent))
	if this.signHandleMethod == "ToUpper" {
		sign = strings.ToUpper(sign)
	}

	signState := sign == this.inputSign

	format := "signContent:%s, sign:%s, inputSign:%s"
	signMsg := fmt.Sprintf(format, this.signContent, sign, this.inputSign)

	if err = this.Base.CheckSign(signState, this.signMethod, signMsg); err != nil {
		return
	}

	return
}

type Rsa struct {
	Crypto
	signMode      int //0(default):publicKey 1:privateKey
	rsaPublicKey  *rsa.PublicKey
	rsaPrivateKey *rsa.PrivateKey
}

func (this *Rsa) CheckSign(params ...interface{}) (err error) {
	this.signMethod = "rsa"
	rsaFuncName := ""
	signState := false
	signMsg := ""
	payKey := this.channelParams["_payKey"]
	hashType := crypto.SHA1

	switch this.signMode {
	case 0:
		if this.rsaPublicKey, err = tool.ParsePKIXPublicKeyWithStr(payKey); err != nil {
			this.lastError = err_parseRsaPublicKey
			rsaFuncName = "parsePKIXPublicKey"
		} else if err = tool.RsaVerifyPKCS1v15(this.rsaPublicKey, hashType, this.signContent, this.inputSign); err != nil {
			this.lastError = err_parseRsaPublicKey
			rsaFuncName = "verifyPKCS1v15"
		}
		signState = err == nil
		signMsg = fmt.Sprintf("inputSign:%s, signContent:%s", this.inputSign, this.signContent)
	case 1:
		sign := ""
		if this.rsaPrivateKey, err = tool.ParsePkCS8PrivateKeyWithStr(payKey); err != nil {
			this.lastError = err_parseRsaPrivateKey
			rsaFuncName = "parsePkCS8PrivateKey"
		} else if sign, err = tool.RsaPKCS1V15Sign(this.rsaPrivateKey, hashType, this.signContent); err != nil {
			rsaFuncName = "rsaPKCS1V15Sign"
		}
		signState = sign == this.inputSign
		format := "inputSign:%s, signContent:%s, sign:%s"
		signMsg = fmt.Sprintf(format, this.inputSign, this.signContent, sign)
	}

	if err != nil {
		format := "CheckSign: %s is error, signMethod:%s, err:%+v, channelId:%d, productId:%d"
		msg := fmt.Sprintf(format, rsaFuncName, this.signMethod, err, this.channelId, this.productId)
		beego.Error(msg)
		return
	}

	if err = this.Base.CheckSign(signState, this.signMethod, signMsg); err != nil {
		return
	}
	return
}
