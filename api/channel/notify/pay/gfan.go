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
	"u9/tool"
)

var gfanUrlKeys []string = []string{"time", "sign"}

const (
	err_gfanParseUId  = 10601
	err_gfanParseBody = 10602
)

type GFan struct {
	Base
	uid       string
	timeStamp string
	sign      string
	ctx       *context.Context
}

func NewGFan(channelId, productId int, urlParams *url.Values, ctx *context.Context) *GFan {
	ret := new(GFan)
	ret.Init(channelId, productId, urlParams, ctx)
	return ret
}

func (this *GFan) Init(channelId, productId int, urlParams *url.Values, ctx *context.Context) {
	this.Base.Init(channelId, productId, urlParams, &gfanUrlKeys)
	this.ctx = ctx
}

func (this *GFan) parseUId() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_gfanParseUId
			beego.Trace(err)
		}
	}()
	this.uid, err = this.getPackageParam("GFAN_UID")
	return
}

func (this *GFan) parseUrlParam() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_parseUrlParam
			beego.Trace(err)
		}
	}()

	this.timeStamp = this.urlParams.Get("time")
	this.sign = this.urlParams.Get("sign")
	return
}

func (this *GFan) parseBody() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_gfanParseBody
			beego.Trace(err)
		}
	}()

	type Response struct {
		OrderId    string `xml:"order_id"`
		Cost       int    `xml:"cost"`
		Appkey     uint64 `xml:"appkey"`
		CreateTime uint64 `xml:"create_time"`
	}

	var buffer bytes.Buffer
	if _, err = io.Copy(&buffer, this.ctx.Request.Body); err != nil {
		return
	}
	content := string(buffer.Bytes())

	var response Response
	if err = xml.Unmarshal([]byte(content), &response); err != nil {
		beego.Trace(content)
		return
	}

	this.orderId = response.OrderId
	this.payAmount = response.Cost * 10 //1元 = 10机锋券
	return
}

func (this *GFan) ParseParam() (err error) {
	if err = this.parseUId(); err != nil {
		return
	}
	if err = this.parseUrlParam(); err != nil {
		return
	}
	if err = this.parseBody(); err != nil {
		return
	}
	if err = this.Base.ParseParam(); err != nil {
		return
	}
	this.channelUserId = this.loginRequest.ChannelUserid
	this.channelOrderId = this.orderRequest.ChannelOrderId
	return
}

func (this *GFan) CheckSign() (err error) {
	defer func() {
		if err != nil {
			this.callbackRet = err_checkSign
			beego.Trace(err)
		}
	}()

	content := this.uid + this.timeStamp
	if sign := tool.Md5([]byte(content)); sign != this.sign {
		msg := fmt.Sprintf("Sign is invalid, content:%s, sign:%s, urlSign%s", content, sign, this.sign)
		err = errors.New(msg)
		return
	}
	return
}

func (this *GFan) GetResult() (ret string) {
	if this.callbackRet == err_noerror {
		ret = "<response><ErrorCode>1</ErrorCode><ErrorDesc>success</ErrorDesc></response>"
	} else {
		beego.Trace(this.callbackRet)
		ret = "<response><ErrorCode>0</ErrorCode><ErrorDesc>fail</ErrorDesc></response>"
	}
	return
}
