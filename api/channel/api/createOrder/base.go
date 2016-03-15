package createOrder

import (
	"u9/api/common"
	"u9/models"
)

type CreateOrder interface {
	InitParam() (err error)
	Response() (err error)
	ParseChannelRet() (err error)
	GetResult() (ret string)
	GetChannelOrderId() (ret string)
}

type Cr struct {
	common.Request
	lr             *models.LoginRequest
	orderId        string
	host           string
	urlJsonParam   string
	ChannelOrderId string
}

func (this *Cr) Init(lr *models.LoginRequest, orderId, host, urlJsonParam string) {
	this.Request.Init()
	this.lr = lr
	this.orderId = orderId
	this.host = host
	this.urlJsonParam = urlJsonParam
}

func (this *Cr) GetChannelOrderId() (ret string) {
	return this.ChannelOrderId
}
