package loginRequest

import (
	// "encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"u9/models"
	// "u9/tool"
)

//木蚂蚁

type MuMaYi struct {
	Lr
}

func LrNewMuMaYi(mlr *models.LoginRequest, args *map[string]interface{}) *MuMaYi {
	ret := new(MuMaYi)
	ret.Init(mlr, args)
	return ret
}

func (this *MuMaYi) Init(mlr *models.LoginRequest, args *map[string]interface{}) {
	this.Lr.Init(mlr)

	this.Method = "POST"
	format := "http://pay.mumayi.com/user/index/validation?uid=%s&token=%s"
	this.Url = fmt.Sprintf(format, this.mlr.ChannelUserid, this.mlr.Token)
	beego.Trace(this.Url)
}

func (this *MuMaYi) ParseChannelRet() (err error) {
	beego.Trace(this.Result)
	// return json.Unmarshal([]byte(this.Result), &this.channelRet)
	return
}

func (this *MuMaYi) CheckChannelRet() bool {
	// string(this.Request) == "success"
	return true
}
