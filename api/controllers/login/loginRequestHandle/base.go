package loginRequestHandle

import (
	"u9/api/common"
)

type LRHandle interface {
	Handle() (ret string, err error)
	Init(param *Param) (err error)
	GetToken() (ret string)
}

type LRH struct {
	common.Request
	param *Param
}

func (this *LRH) Init(param *Param) {
	this.Request.Init()
	this.param = param
}
