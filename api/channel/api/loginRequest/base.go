package loginRequest

import (
	"u9/api/common"
	"u9/models"
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
	ret common.BasicRet
	mlr *models.LoginRequest
}

func (this *Lr) Init(mlr *models.LoginRequest) {
	this.Request.Init()
	this.ret.Init()
	this.mlr = mlr
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
