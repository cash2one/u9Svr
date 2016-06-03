package channelPayNotify

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"net/url"
	"strconv"
	"u9/tool"
)

type Baidu struct {
	Base
	channelRet baiduChannelRet
}

type baiduResData struct {
	UID             string `json:"UID"`
	MerchandiseName string `json:"MerchandiseName"`
	OrderMoney      string `json:"OrderMoney"`
	StartDateTime   string `json:"StartDateTime"`
	BankDateTime    string `json:"BankDateTime"`
	OrderStatus     int    `json:"OrderStatus"`
	StatusMsg       string `json:"StatusMsg"`
	ExtInfo         string `json:"ExtInfo"`
	VoucherMoney    int    `json:"VoucherMoney"`
}

type baiduChannelRet struct {
	AppID      string `json:"AppID"`
	ResultCode string `json:"ResultCode"`
	ResultMsg  string `json:"ResultMsg"`
	Sign       string `json:"Sign"`
}

func NewBaidu(channelId, productId int, urlParams *url.Values, ctx *context.Context) *Baidu {
	ret := new(Baidu)
	ret.Init(channelId, productId, urlParams, ctx)
	return ret
}

func (this *Baidu) Init(channelId, productId int, urlParams *url.Values, ctx *context.Context) {
	this.Base.InitWithCtx(channelId, productId, urlParams, &emptyUrlKeys, ctx)
	this.data = new(baiduResData)
}

func (this *Baidu) parseUrlParam() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_htcParseBody
			beego.Error(this.urlParams)
			beego.Error(err)
		}
	}()

	content := this.urlParams.Get("Content")
	var enByte []byte
	if enByte, err = base64.StdEncoding.DecodeString(content); err != nil {
		return
	}
	if err = json.Unmarshal(enByte, &this.data); err != nil {
		return
	}

	this.orderId = this.urlParams.Get("CooperatorOrderSerial")
	this.channelOrderId = this.urlParams.Get("OrderSerial")

	data := this.data.(*baiduResData)
	this.channelUserId = data.UID
	payAmount := 0.0
	if payAmount, err = strconv.ParseFloat(data.OrderMoney, 64); err != nil {
		return err
	} else {
		this.payAmount = int(payAmount * 100)
	}

	return
}

func (this *Baidu) ParseParam() (err error) {
	if err = this.parseUrlParam(); err != nil {
		return
	}
	if err = this.parseChannelGameID("BAIDU_APPID"); err != nil {
		return
	}
	if err = this.parseChannelPayKey("BAIDU_APPSECRET"); err != nil {
		return
	}
	if err = this.Base.ParseParam(); err != nil {
		return
	}
	return
}

func (this *Baidu) CheckSign() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_checkSign
			beego.Error(err)
		}
	}()

	content := this.urlParams.Get("AppID") + this.urlParams.Get("OrderSerial") +
		this.urlParams.Get("CooperatorOrderSerial") +
		this.urlParams.Get("Content") + this.channelPayKey

	urlSign := this.urlParams.Get("Sign")
	if sign := tool.Md5([]byte(content)); sign != urlSign {
		msg := fmt.Sprintf("sign(%s) isn't equal urlSign(%s)", sign, urlSign)
		err = errors.New(msg)
		return
	}
	return
}

func (this *Baidu) ParseChannelRet() (err error) {
	if err = this.Base.ParseChannelRet(); err != nil {
		return
	}

	if this.urlParams.Get("AppID") != this.channelGameId {
		this.callbackRet = err_parseChannelGameId
		beego.Error("BAIDU_APPID is invalid.")
		return
	}

	data := this.data.(*baiduResData)
	if data.OrderStatus != 1 {
		this.callbackRet = err_callbackFail
		beego.Error("transData.result is equal 1.")
		return
	}
	return
}

func (this *Baidu) GetResult() (ret string) {
	beego.Trace("callbackRet:" + strconv.Itoa(this.callbackRet))
	this.channelRet.AppID = this.urlParams.Get("AppID")

	if this.callbackRet == err_noerror {
		this.channelRet.ResultCode = "1"
		this.channelRet.ResultMsg = "成功"
	} else if this.callbackRet == err_checkSign {
		this.channelRet.ResultCode = "1001"
		this.channelRet.ResultMsg = "Sign无效"
	} else {
		this.channelRet.ResultCode = "0"
		this.channelRet.ResultMsg = "其它"
	}

	content := this.channelRet.AppID + this.channelRet.ResultCode + this.channelPayKey
	this.channelRet.Sign = tool.Md5([]byte(content))

	data, _ := json.Marshal(this.channelRet)
	ret = string(data)
	beego.Trace(ret)
	return
}
