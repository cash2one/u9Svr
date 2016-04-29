package channelPayNotify

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"net/url"
	"strconv"
	"u9/tool"
)

var qmyxUrlKeys []string = []string{}

const (
	err_qmyxParseGameKey  = 14301
	err_qmyxParseBody     = 14302
	err_qmyxResultFailure = 14303
)

type Qmyx struct {
	Base
	gameId   string
	gameKey  string
	qmyxData QmyxData
}

type QmyxData struct {
	ChannelOrderId string `json:"pay_order_code"`
	GameId         string `json:"game_id"`
	ServerId       string `json:"server_id"`
	UserKey        string `json:"user_key"`
	Amount         int    `json:"amount"`
	Ext            string `json:"ext"`
}

func NewQmyx(channelId, productId int, urlParams *url.Values) *Qmyx {
	ret := new(Qmyx)
	ret.Init(channelId, productId, urlParams)
	return ret
}

func (this *Qmyx) Init(channelId, productId int, urlParams *url.Values) {
	this.Base.Init(channelId, productId, urlParams, &qmyxUrlKeys)
}

func (this *Qmyx) parseGameID() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_parseProductKey
			beego.Trace(err)
		}
	}()
	this.gameId, err = this.getPackageParam("GAME_ID")
	return
}

func (this *Qmyx) parseGameKey() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_qmyxParseGameKey
			beego.Trace(err)
		}
	}()
	this.gameKey, err = this.getPackageParam("GAME_KEY")
	return
}

func (this *Qmyx) parseUrlParam() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_qmyxParseBody
			beego.Error(err)
			beego.Error(this.urlParams)
		}
	}()

	data := this.urlParams.Get("data")
	if err = json.Unmarshal([]byte(data), &this.qmyxData); err != nil {
		return err
	}

	this.orderId = this.qmyxData.Ext
	this.channelOrderId = this.qmyxData.ChannelOrderId
	this.channelUserId = this.qmyxData.UserKey
	this.payAmount = this.qmyxData.Amount
	return
}

func (this *Qmyx) ParseChannelRet() (err error) {
	if this.orderId != this.orderRequest.OrderId {
		this.callbackRet = err_orderIsNotExist
		return
	}

	if this.orderRequest.ReqAmount != this.payAmount {
		this.callbackRet = err_payAmountError
		return
	}

	if this.qmyxData.GameId != this.gameId {
		this.callbackRet = err_parseProductKey
		return
	}
	return
}

func (this *Qmyx) ParseParam() (err error) {
	if err = this.parseUrlParam(); err != nil {
		return
	}
	if err = this.parseGameID(); err != nil {
		return
	}
	if err = this.parseGameKey(); err != nil {
		return
	}
	if err = this.Base.ParseParam(); err != nil {
		return
	}
	return
}

func (this *Qmyx) CheckSign() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_checkSign
			beego.Trace(err)
		}
	}()

	content := this.qmyxData.ChannelOrderId + this.qmyxData.GameId +
		this.qmyxData.ServerId + this.qmyxData.UserKey +
		strconv.Itoa(this.qmyxData.Amount) + this.qmyxData.Ext + this.gameKey

	urlSign := this.urlParams.Get("sign")
	if sign := tool.Md5([]byte(content)); sign != urlSign {
		msg := fmt.Sprintf("Sign is invalid, content:%s, sign:%s, urlSign:%s", content, sign, urlSign)
		err = errors.New(msg)
		return
	}
	return
}

func (this *Qmyx) GetResult() (ret string) {
	beego.Trace(this.callbackRet)
	if this.callbackRet == err_noerror {
		ret = `success`
	} else {
		ret = `failure`
	}
	return
}
