package channelPayNotify

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"net/url"
	"sort"
	"strconv"
	"u9/tool"
)

type AnFeng struct {
	Base
}

type anFengResData struct {
	Uid        string `json:"uid"`
	Ucid       string `json:"ucid"`
	Body       string `json:"body"`
	Fee        string `json:"fee"`
	Subject    string `json:"subject"`
	Vid        string `json:"vid"`
	Sn         string `json:"sn"`
	Vorderid   string `json:"vorderid"`
	CreateTime string `json:"createTime"`
	Sign       string `json:"sign"`
}

func NewAnFeng(channelId, productId int, urlParams *url.Values, ctx *context.Context) *AnFeng {
	ret := new(AnFeng)
	ret.Init(channelId, productId, urlParams, ctx)
	return ret
}

func (this *AnFeng) Init(channelId, productId int, urlParams *url.Values, ctx *context.Context) {
	this.Base.InitWithCtx(channelId, productId, urlParams, &emptyUrlKeys, ctx)
}

func (this *AnFeng) parseUrlParam() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_htcParseBody
			beego.Error(this.urlParams)
			beego.Error(err)
		}
	}()

	this.orderId = this.urlParams.Get("vorderid")
	this.channelOrderId = this.urlParams.Get("sn")
	this.channelUserId = this.urlParams.Get("ucid")
	payAmount := 0.0
	if payAmount, err = strconv.ParseFloat(this.urlParams.Get("fee"), 64); err != nil {
		return err
	} else {
		this.payAmount = int(payAmount * 100)
	}

	return
}

func (this *AnFeng) ParseParam() (err error) {
	if err = this.parseUrlParam(); err != nil {
		return
	}
	if err = this.parseChannelGameID("productid"); err != nil {
		return
	}
	if err = this.parseChannelPayKey("ANFENG_PAYKEY"); err != nil {
		return
	}
	if err = this.Base.ParseParam(); err != nil {
		return
	}
	return
}

func (this *AnFeng) CheckSign() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_checkSign
			beego.Error(err)
		}
	}()

	excludeItems := []string{"sign"}
	sorter := tool.NewUrlValuesSorter(this.urlParams, &excludeItems)
	sort.Sort(sorter)
	content := sorter.DefaultBody() + "&signKey=" + this.channelPayKey
	beego.Trace(content)
	urlSign := this.urlParams.Get("sign")
	if sign := tool.Md5([]byte(content)); sign != urlSign {
		msg := fmt.Sprintf("Sign is invalid, content:%s, urlSign:%s", content, urlSign)
		err = errors.New(msg)
		return
	}

	return
}

func (this *AnFeng) ParseChannelRet() (err error) {
	if err = this.Base.ParseChannelRet(); err != nil {
		return
	}

	if this.urlParams.Get("vid") != this.channelGameId {
		this.callbackRet = err_parseChannelGameId
		beego.Error("productid is invalid.")
		return
	}
	return
}

func (this *AnFeng) GetResult() (ret string) {
	beego.Trace("callbackRet:" + strconv.Itoa(this.callbackRet))
	if this.callbackRet == err_noerror {
		ret = "SUCCESS"
	} else {
		ret = "FAILURE"
	}
	return
}
