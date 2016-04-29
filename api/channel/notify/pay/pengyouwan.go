package channelPayNotify

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"io"
	"net/url"
	"strconv"
	"u9/tool"
)

var pengyouwanUrlKeys []string = []string{}

const (
	err_pengyouwanParsePayKey   = 14201
	err_pengyouwanParseBody     = 14202
	err_pengyouwanResultFailure = 14203
)

type Pengyouwan struct {
	Base
	gameKey        string
	payKey         string
	pengyouwanData PengyouwanData
	ctx            *context.Context
}

type PengyouwanData struct {
	Ver       string `json:"ver"`
	Tid       string `json:"tid"`
	Sign      string `json:"sign"`
	Gamekey   string `json:"gamekey"`
	Channel   string `json:"channel"`
	CPOrderid string `json:"cp_orderid"`
	ChOrderid string `json:"ch_orderid"`
	Amount    string `json:"amount"`
	CpParam   string `json:"cp_param"`
	Ack       int    `json:"ack"`
	Msg       string `json:"msg"`
}

func NewPengyouwan(channelId, productId int, urlParams *url.Values, ctx *context.Context) *Pengyouwan {
	ret := new(Pengyouwan)
	ret.Init(channelId, productId, urlParams, ctx)
	return ret
}

func (this *Pengyouwan) Init(channelId, productId int, urlParams *url.Values, ctx *context.Context) {
	this.Base.Init(channelId, productId, urlParams, &pengyouwanUrlKeys)
	this.ctx = ctx
}

func (this *Pengyouwan) parseGameKey() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_pengyouwanParsePayKey
			beego.Trace(err)
		}
	}()
	this.gameKey, err = this.getPackageParam("PYW_GAME_KEY")
	return
}

func (this *Pengyouwan) parsePayKey() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_pengyouwanParsePayKey
			beego.Trace(err)
		}
	}()
	this.payKey, err = this.getPackageParam("PYW_PAY_KEY")
	return
}

func (this *Pengyouwan) parseBody() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_pengyouwanParseBody
			beego.Error(err)
		}
	}()

	var buffer bytes.Buffer
	if _, err = io.Copy(&buffer, this.ctx.Request.Body); err != nil {
		return
	}

	body, _ := url.QueryUnescape(string(buffer.Bytes()))
	beego.Trace(body)

	if err = json.Unmarshal([]byte(body), &this.pengyouwanData); err != nil {
		return err
	}

	this.orderId = this.pengyouwanData.CPOrderid
	this.channelOrderId = this.pengyouwanData.ChOrderid

	payAmount := 0.0
	if payAmount, err = strconv.ParseFloat(this.pengyouwanData.Amount, 64); err != nil {
		return err
	} else {
		this.payAmount = int(payAmount * 100)
	}
	return
}

func (this *Pengyouwan) ParseChannelRet() (err error) {
	if this.orderId != this.orderRequest.OrderId {
		this.callbackRet = err_orderIsNotExist
		return
	}

	if this.orderRequest.ReqAmount != this.payAmount {
		this.callbackRet = err_payAmountError
		return
	}

	if this.pengyouwanData.Gamekey != this.gameKey {
		this.callbackRet = err_parseProductKey
		return
	}

	if this.pengyouwanData.Ack != 200 {
		this.callbackRet = err_pengyouwanResultFailure
		return
	}
	return
}

func (this *Pengyouwan) ParseParam() (err error) {
	if err = this.parseBody(); err != nil {
		return
	}
	if err = this.parseGameKey(); err != nil {
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

func (this *Pengyouwan) CheckSign() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_checkSign
			beego.Trace(err)
		}
	}()

	content := this.payKey + this.pengyouwanData.CPOrderid + this.pengyouwanData.ChOrderid + this.pengyouwanData.Amount

	urlSign := this.pengyouwanData.Sign
	if sign := tool.Md5([]byte(content)); sign != urlSign {
		msg := fmt.Sprintf("Sign is invalid, content:%s, sign:%s, urlSign:%s", content, sign, urlSign)
		err = errors.New(msg)
		return
	}
	return
}

func (this *Pengyouwan) GetResult() (ret string) {
	beego.Trace(this.callbackRet)
	if this.callbackRet == err_noerror {
		ret = `{"ack":200,"ok"}`
	} else {
		ret = `{"ack":500,"fail"}`
	}
	return
}
