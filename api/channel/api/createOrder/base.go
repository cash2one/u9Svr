package createOrder

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"u9/api/common"
	"u9/models"
)

type CreateOrder interface {
	Prepare(lr *models.LoginRequest, orderId, extParamStr string,
		channelParams *map[string]interface{}, ctx *context.Context) (err error)
	InitParam() (err error)
	GetResponse() (err error)
	ParseChannelRet() (err error)
	GetResult() (ret string)
	GetChannelOrderId() (ret string)
}

type Cr struct {
	common.Request
	ctx            *context.Context
	lr             *models.LoginRequest
	channelParams  *map[string]interface{}
	extParam       interface{}
	channelRet     interface{}
	orderId        string
	channelOrderId string
	extParamStr    string
	appId          string
	appKey         string
	payKey         string
}

func (this *Cr) Initial(
	lr *models.LoginRequest,
	orderId string,
	channelRet interface{},
	extParam interface{},
	extParamStr string,
	channelParams *map[string]interface{},
	ctx *context.Context) (err error) {

	this.Request.Init()
	this.lr = lr
	this.orderId = orderId
	this.ctx = ctx
	this.channelParams = channelParams

	this.channelRet = channelRet
	this.extParamStr = extParamStr
	this.extParam = extParam

	if this.extParam != nil {
		if err = json.Unmarshal([]byte(this.extParamStr), &this.extParam); err != nil {
			format := "initial: err:%v "
			msg := fmt.Sprintf(format, err) + this.Dump()
			err = errors.New(msg)
			beego.Error(err)
			return
		}
	}

	return
}

func (this *Cr) Dump() (ret string) {
	format := "this:%+v"
	ret = fmt.Sprintf(format, this)
	return
}

func (this *Cr) parseAppKey(key string) {
	this.appKey = (*this.channelParams)[key].(string)
}

func (this *Cr) parseAppId(key string) {
	this.appId = (*this.channelParams)[key].(string)
}

func (this *Cr) parsePayKey(key string) {
	this.payKey = (*this.channelParams)[key].(string)
}

func (this *Cr) GetChannelOrderId() (ret string) {
	return this.channelOrderId
}

func (this *Cr) ParseChannelRet() (err error) {
	//beego.Trace("cr-parseChannelRet: Result:" + this.Result)
	if err = json.Unmarshal([]byte(this.Result), &this.channelRet); err != nil {
		format := "cr-parseChannelRet: err:%v, result:%s, channelRet:%+v"
		msg := fmt.Sprintf(format, err, this.Result, this.channelRet)
		err = errors.New(msg)
		beego.Error(err)
		return
	}
	return nil
}
