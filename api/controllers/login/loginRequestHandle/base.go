package loginRequestHandle

import (
	"u9/api/common"
)

type LRHandle interface {
	Init(param *Param) (err error)
	Handle() (ret string, err error)
	//GetChannelResult() (ret interface{})
}

type LRH struct {
	common.Request
	param *Param
}

func (this *LRH) Init(param *Param) {
	this.Request.Init()
	this.param = param
}
