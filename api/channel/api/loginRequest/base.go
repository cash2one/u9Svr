package loginRequest

import (
	"u9/api/common"
)

type LoginRequest interface {
	InitParam()
	Response() (err error)
	ParseChannelRet() (err error)
	CheckChannelRet() bool
	SetCode(code int) *common.BasicRet
}

type Lr struct {
	common.Request
	ret           common.BasicRet
	channelUserId string
	token         string
}

func (this *Lr) Init(channelUserId, token string) {
	this.Request.Init()
	this.ret.Init()
	this.channelUserId = channelUserId
	this.token = token
}

func (this *Lr) SetCode(code int) *common.BasicRet {
	this.ret.SetCode(code)
	return &this.ret
}

func (this *Lr) ParseChannelRet() (err error) {
	return nil
}

func (this *Lr) CheckChannelRet() bool {
	return false
}
