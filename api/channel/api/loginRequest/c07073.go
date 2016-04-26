package loginRequest

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"u9/models"
	"u9/tool"
)

//07073

type C07073ChannelRet struct {
	State   int 	 `json:"state"`
	Msg     string   `json:"msg"`
	Data    C07073DataJson   `json:"data"`
}

type C07073DataJson struct{
	Username string `json:"username"`
	Uid string `json:"uid"`
}

type C07073 struct {
	Lr
	channelRet C07073ChannelRet
}

func LrNewC07073(mlr *models.LoginRequest, args *map[string]interface{}) *C07073 {
	ret := new(C07073)
	ret.Init(mlr, args)
	return ret
}

func (this *C07073) Init(mlr *models.LoginRequest, args *map[string]interface{}) {
	this.Lr.Init(mlr)
	pid := (*args)["C07073_PID"].(string)
	key := (*args)["C07073_SECRET_KEY"].(string)
	sign := tool.Md5([]byte("pid=" + pid +"&token="+  this.mlr.Token + "&username=" +
		this.mlr.ChannelUserid + key))
	format := "http://sdk.07073sy.com/index.php/User/v4?username=%s&token=%s&pid=%s&sign=%s"
	this.Url = fmt.Sprintf(format, mlr.ChannelUserid ,this.mlr.Token ,pid, sign)
 	beego.Trace(this.Url)
}

func (this *C07073) ParseChannelRet() (err error) {
	beego.Trace(this.Result)
	return json.Unmarshal([]byte(this.Result), &this.channelRet)
}

func (this *C07073) CheckChannelRet() bool {
	return this.channelRet.State == 1
}
