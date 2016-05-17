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

type Qmyx struct {
	Base
	gameId string
}

type qmyxData struct {
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
	this.Base.Init(channelId, productId, urlParams, &emptyUrlKeys)
	this.data = new(qmyxData)
}

func (this *Qmyx) parseUrlParam() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_parseUrlParam
			beego.Error(err)
			beego.Error(this.urlParams)
		}
	}()

	param := this.urlParams.Get("data")
	if err = json.Unmarshal([]byte(param), &this.data); err != nil {
		return err
	}

	data := this.data.(*qmyxData)

	this.orderId = data.Ext
	this.channelOrderId = data.ChannelOrderId
	this.channelUserId = data.UserKey
	this.payAmount = data.Amount
	return
}

func (this *Qmyx) ParseChannelRet() (err error) {
	if err = this.Base.ParseChannelRet(); err != nil {
		this.callbackRet = err_parseChannelRet
		return
	}

	data := this.data.(*qmyxData)
	if data.GameId != this.channelGameId {
		this.callbackRet = err_parseProductKey
		return
	}
	return
}

func (this *Qmyx) ParseParam() (err error) {
	if err = this.parseUrlParam(); err != nil {
		return
	}
	if err = this.parseChannelGameID("GAME_ID"); err != nil {
		return
	}
	if err = this.parseChannelGameKey("GAME_KEY"); err != nil {
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
			beego.Error(err)
		}
	}()

	data := this.data.(*qmyxData)
	content := data.ChannelOrderId + data.GameId + data.ServerId + data.UserKey +
		strconv.Itoa(data.Amount) + data.Ext + this.channelGameKey

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
