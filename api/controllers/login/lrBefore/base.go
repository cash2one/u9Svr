package lrBefore

import (
	"errors"
	"fmt"
	"u9/api/common"
	"u9/models"
)

type Handle interface {
	Init(param *Param) (err error)
	Exec() (ret string, err error)
}

type base struct {
	common.Request
	param         *Param
	channelParams map[string]string
	channelRet    interface{}
}

func (this *base) Init(param *Param) {
	this.Request.Init()
	this.param = param

	this.initChannelParam()
}

func (this *base) Dump() (ret string) {
	format := "param:%+v, url:%s, result:%s, channelRet:%+v"
	ret = fmt.Sprintf(format, this.param, this.Url, this.Result, this.channelRet)
	return
}

func (this *base) initChannelParam() (err error) {
	this.channelParams = map[string]string{}

	channelId := this.param.ChannelId
	productId := this.param.ProductId
	if this.channelParams, err = models.GetChannelXmlParam(channelId, productId); err != nil {
		format := "getChannelXmlParam is err:%v"
		msg := fmt.Sprintf(format, err)
		err = errors.New(msg)
		return
	}
	return
}

func (this *base) getChannelParam(paramKey string) (ret string, err error) {
	ok := false
	if ret, ok = this.channelParams[paramKey]; !ok {
		msg := fmt.Sprintf("packageParam %s isn't exist", paramKey)
		err = errors.New(msg)
		return
	}
	return
}
