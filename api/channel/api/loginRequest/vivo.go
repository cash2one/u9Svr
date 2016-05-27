package loginRequest

import (
	"encoding/json"
	//"fmt"
	"github.com/astaxie/beego"
	"u9/models"
	//"u9/tool"
)

type vivoChannelRet struct {
	Stat  string `json:"stat"`
	Msg   string `json:"msg"`
	Uid   string `json:"uid"`
	Email string `json:"email"`
}

type Vivo struct {
	Lr
	channelRet vivoChannelRet
	args       *map[string]interface{}
}

func LrNewVivo(mlr *models.LoginRequest, args *map[string]interface{}) *Vivo {
	ret := new(Vivo)
	ret.Init(mlr, args)
	return ret
}

func (this *Vivo) Init(mlr *models.LoginRequest, args *map[string]interface{}) {
	this.Lr.Init(mlr)
	this.Method = "POST"
	this.IsHttps = true
	this.args = args

	this.Url = "https://usrsys.vivo.com.cn/auth/user/info"
	beego.Trace(this.Url)
}

func (this *Vivo) InitParam() (err error) {
	if err = this.Lr.InitParam(); err != nil {
		return
	}
	this.Req.Param("access_token", this.mlr.Token)
	return
}

func (this *Vivo) ParseChannelRet() (err error) {
	beego.Trace(this.Result)
	return json.Unmarshal([]byte(this.Result), &this.channelRet)
}

func (this *Vivo) CheckChannelRet() bool {
	return this.channelRet.Uid != ""
}
