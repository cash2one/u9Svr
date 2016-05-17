package channelPayNotify

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"io"
	"net/url"
	"strconv"
	"u9/tool"
)

type Jiuyou struct {
	Base
}

type jiuyouResData_data struct {
	OrderId      string `json:"orderId"`
	GameId       int    `json:"gameId"`
	AccountId    string `json:"accountId"`
	Creator      string `json:"creator"`
	PayWay       int    `json:"payWay"`
	Amount       string `json:"amount"`
	CallbackInfo string `json:"callbackInfo"`
	OrderStatus  string `json:"orderStatus"`
	FailedDesc   string `json:"failedDesc"`
	CpOrderId    string `json:"cpOrderId"`
}

type jiuyouResData struct {
	Ver  string             `json:"ver"`
	Data jiuyouResData_data `json:"data"`
	Sign string             `json:"sign"`
}

func NewJiuyou(channelId, productId int, urlParams *url.Values, ctx *context.Context) *Jiuyou {
	ret := new(Jiuyou)
	ret.Init(channelId, productId, urlParams, ctx)
	return ret
}

func (this *Jiuyou) Init(channelId, productId int, urlParams *url.Values, ctx *context.Context) {
	this.Base.InitWithCtx(channelId, productId, urlParams, &emptyUrlKeys, ctx)
	this.data = new(jiuyouResData)
}

func (this *Jiuyou) parseBody() (err error) {
	var body string
	defer func() {
		if err != nil {
			this.callbackRet = err_parseBody
			beego.Error(err)
			beego.Error("body:" + body)
		}
	}()

	var buffer bytes.Buffer
	if _, err = io.Copy(&buffer, this.ctx.Request.Body); err != nil {
		return
	}

	body = string(buffer.Bytes())
	if err = xml.Unmarshal([]byte(body), &this.data); err != nil {
		return
	}

	data := this.data.(*jiuyouResData)
	beego.Trace(data)
	//if this.orderId, err = strconv.Atoi(data.Data.CpOrderId); err != nil {
	//	return
	//}
	this.orderId = data.Data.CpOrderId
	this.channelOrderId = data.Data.OrderId
	this.channelUserId = data.Data.AccountId

	payAmount := 0.0
	if payAmount, err = strconv.ParseFloat(data.Data.Amount, 64); err != nil {
		return err
	} else {
		this.payAmount = int(payAmount * 100)
	}

	return
}

func (this *Jiuyou) ParseParam() (err error) {
	if err = this.parseBody(); err != nil {
		return
	}
	if err = this.parseChannelPayKey("UC_APPKEY"); err != nil {
		return
	}
	if err = this.Base.ParseParam(); err != nil {
		return
	}
	return
}

func (this *Jiuyou) CheckSign() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_checkSign
			beego.Error(err)
		}
	}()

	data := this.data.(*jiuyouResData)
	beego.Trace(data)

	format := `accountId=%samount=%scallbackInfo=%scpOrderId=%screator=%sfailedDesc=%sgameId=%sorderId=%sorderStatus=%spayWay=%s%s`

	content := fmt.Sprintf(format,
		data.Data.AccountId, data.Data.AccountId, data.Data.CallbackInfo, data.Data.CpOrderId,
		data.Data.Creator, data.Data.FailedDesc, data.Data.GameId, data.Data.OrderId, data.Data.OrderStatus,
		data.Data.PayWay, this.channelPayKey)

	//content := this.data.ChannelOrderId + this.data.GameId +
	//	this.data.ServerId + this.data.UserKey +
	//	strconv.Itoa(this.data.Amount) + this.data.Ext + this.gameKey

	urlSign := data.Sign
	if sign := tool.Md5([]byte(content)); sign != urlSign {
		msg := fmt.Sprintf("Sign is invalid, content:%s, sign:%s, urlSign:%s", content, sign, urlSign)
		err = errors.New(msg)
		return
	}
	return
}

func (this *Jiuyou) GetResult() (ret string) {
	beego.Trace("callbackRet:" + strconv.Itoa(this.callbackRet))
	if this.callbackRet == err_noerror {
		ret = `SUCCESS`
	} else {
		ret = `FAILURE`
	}
	return
}
