package channelPayNotify

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"bytes"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"strconv"
	"io"
	"u9/tool"
	"net/url"
)

var ttUrlKeys []string = []string{}

const (
	err_ttParsePayKey      = 12701
	err_ttResultFailure    = 12702
	err_ttInitRsaPublicKey = 12703
	err_ttParseBody        = 12704
)

//tt
type TT struct {
	Base
	payKey     string
	tt_Sign    string
	tt_result TT_Result
	tt_contentBody string
	ctx        *context.Context
}

type TT_Result struct {
	Uid   		 int      `json:"uid"`
	GameId     	 int 	   `json:"gameId"`
	SDKOrderId   string     `json:"sdkOrderId"`
	CpOrderId    string 	`json:"cpOrderId"`
	PayFee     	 string 	`json:"payFee"`
	PayResult 	 string 	`json:"payResult"`
	PayDate 	 string 	`json:"payDate"`
	ExInfo   	 string 	`json:"exInfo"`
}

var (
	ttRsaPublicKey *rsa.PublicKey
)

func NewTT(channelId, productId int, urlParams *url.Values, ctx *context.Context) *TT {
	ret := new(TT)
	ret.Init(channelId, productId, urlParams, ctx)
	return ret
}

func (this *TT) Init(channelId, productId int, urlParams *url.Values, ctx *context.Context) {
	this.Base.Init(channelId, productId, urlParams, &ttUrlKeys)
	this.ctx = ctx
}

func (this *TT) parsePayKey() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_ttParsePayKey
			beego.Trace(err)
		}
	}()
	this.payKey, err = this.getPackageParam("TT_SDK_PAYKEY")
	return
}

func (this *TT) CheckUrlParam() (err error) {
	return
}

func (this *TT) parseUrlParam() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_parseUrlParam
			beego.Trace(err)
		}
	}()

	// beego.Trace(this.response.Order)
	beego.Trace(this.tt_result)
	this.orderId = this.tt_result.CpOrderId
	this.channelOrderId = this.tt_result.SDKOrderId
	payAmount := 0.0
	if payAmount, err = strconv.ParseFloat(this.tt_result.PayFee, 64); err != nil {
		return err
	} else {
		this.payAmount = int(payAmount * 100)
	}
	return
}

func (this *TT) ParseChannelRet() (err error) {
	if result := this.tt_result.PayResult; result != "1" {
		this.callbackRet = err_ttResultFailure
	}
	return
}

func (this *TT) parseBody() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_htcParseBody
			beego.Trace(err)
		}
	}()

	var buffer bytes.Buffer
	if _, err = io.Copy(&buffer, this.ctx.Request.Body); err != nil {
		return
	}
	contentBody := string(buffer.Bytes())
	// contentBody = string(url.QueryUnescape(contentBody))
	beego.Trace(contentBody)
	this.tt_contentBody,_ = url.QueryUnescape(contentBody)
	beego.Trace(this.tt_contentBody)
	if err = json.Unmarshal([]byte(this.tt_contentBody), &this.tt_result);err != nil{
		beego.Error(err)
		return err
	}
	
	// beego.Trace(this.ctx.Request)
	// if _, err = io.Copy(&buffer,this.ctx.Request.Head); err != nil {
	// 	return
	// }
	
// string(buffer.Bytes())
	this.tt_Sign = this.ctx.Request.Header.Get("Sign")
	beego.Trace("head OK:"+ this.tt_Sign)
	return
}

func (this *TT) ParseParam() (err error) {
	if err = this.parseBody(); err != nil {
		return
	}
	if err = this.parseUrlParam(); err != nil {
		return
	}
	if err = this.parsePayKey(); err != nil {
		return
	}
	if err = this.Base.ParseParam(); err != nil {
		return
	}
	return
}


func (this *TT) CheckSign() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_checkSign
			beego.Trace(err)
		}
	}()
	// jsonResult,_ := json.Marshal(this.tt_result)
	// format := string(jsonResult)

	content := fmt.Sprintf("%s%s",this.tt_contentBody,this.payKey)
	beego.Trace("content:"+content)
	var result string
	if result,err = tool.TTSign(content);err != nil {
		beego.Error(err)
	}

	// signMd5 := tool.Md5([]byte(content))

	// beego.Trace(signMd5)
	// signMd5 = Substr(signMd5,8,16)
	
	// beego.Trace("md5:",signMd5)

    // sign := base64.StdEncoding.EncodeToString([]byte(signMd5))

    // beego.Trace("sign:",sign)
    // beego.Trace(this.tt_Sign)
    if result != this.tt_Sign{
		msg := fmt.Sprintf("Sign is invalid, sign:%s, urlSign:%s", result, this.tt_Sign)
		err = errors.New(msg)
		return
	}
	return
}

func (this *TT) GetResult() (ret string) {
	if this.callbackRet == err_noerror {
		ret = `{"head":{"result":"0","message":"成功"}}`
	} else {
		ret = `{"head":{"result":"1","message":"失败"}}`
	}
	return
}

func Substr(str string, start, length int) string {
    rs := []rune(str)
    rl := len(rs)
    end := 0
        
    if start < 0 {
        start = rl - 1 + start
    }
    end = start + length
    
    if start > end {
        start, end = end, start
    }
    
    if start < 0 {
        start = 0
    }
    if start > rl {
        start = rl
    }
    if end < 0 {
        end = 0
    }
    if end > rl {
        end = rl
    }
    return string(rs[start:end])
}
