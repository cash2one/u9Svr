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

//朋友玩
type pengyouwanData_cpParam struct {
	ProductId   int    `json:"product_id"`
	OrderId     string `json:"order_id"`
	ProductDesc string `json:"product_desc"`
}

type pengyouwanData struct {
	Ver       string `json:"ver"`
	Tid       string `json:"tid"`
	Sign      string `json:"sign"`
	Gamekey   string `json:"gamekey"`
	Channel   string `json:"channel"`
	CPOrderid string `json:"cp_orderid"`
	ChOrderid string `json:"ch_orderid"`
	Amount    string `json:"amount"`
	CpParam   string `json:"cp_param"`
}

type Pengyouwan struct {
	Base
}

func NewPengyouwan(channelId, productId int, urlParams *url.Values) *Pengyouwan {
	ret := new(Pengyouwan)
	ret.Init(channelId, productId, urlParams)
	return ret
}

func (this *Pengyouwan) Init(channelId, productId int, urlParams *url.Values) {
	this.Base.Init(channelId, productId, urlParams, &emptyUrlKeys)
	this.data = new(pengyouwanData)
}

func (this *Pengyouwan) parseUrlParam() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_parseBody
			beego.Error(err)
			beego.Error(this.urlParams.Encode())
		}
	}()

	content, _ := url.QueryUnescape(this.urlParams.Encode())
	content = content[0 : len(content)-1]
	beego.Trace(content)

	if err = json.Unmarshal([]byte(content), &this.data); err != nil {
		return
	}

	data := this.data.(*pengyouwanData)
	beego.Trace(data)

	this.orderId = data.CPOrderid
	this.channelOrderId = data.ChOrderid

	payAmount := 0.0
	if payAmount, err = strconv.ParseFloat(data.Amount, 64); err != nil {
		return err
	} else {
		this.payAmount = int(payAmount * 100)
	}
	return
}

func (this *Pengyouwan) ParseChannelRet() (err error) {
	if err = this.Base.ParseChannelRet(); err != nil {
		return
	}

	data := this.data.(*pengyouwanData)

	if data.Gamekey != this.channelGameKey {
		this.callbackRet = err_parseChannelGameKey
		return
	}
	return
}

func (this *Pengyouwan) ParseParam() (err error) {
	if err = this.parseUrlParam(); err != nil {
		return
	}
	if err = this.parseChannelGameKey("PYW_GAME_KEY"); err != nil {
		return
	}
	if err = this.parseChannelPayKey("PYW_PAY_KEY"); err != nil {
		return
	}
	if err = this.Base.ParseParam(); err != nil {
		return
	}
	this.channelUserId = this.loginRequest.ChannelUserid
	return
}

func (this *Pengyouwan) CheckSign() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_checkSign
			beego.Trace(err)
		}
	}()

	data := this.data.(*pengyouwanData)
	content := this.channelPayKey + data.CPOrderid + data.ChOrderid + data.Amount

	urlSign := data.Sign
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
