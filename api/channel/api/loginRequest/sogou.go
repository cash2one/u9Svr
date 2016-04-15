package loginRequest

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"u9/tool"
)

//搜狗

type SogouChannelRet struct {
	Result    bool  `json:"result"`
	// Error     ErrorJson   `json:"error"`
}

// type ErrorJson struct {
// 	Code  int    `json:"code"`
// 	Msg  string  `json:"msg"`
// }

type Sogou struct {
	Lr
	args *map[string]interface{}
	channelRet SogouChannelRet
}

func LrNewSogou(channelUserId, token string, args *map[string]interface{}) *Sogou {
	ret := new(Sogou)
	ret.Init(channelUserId, token, args)
	return ret
}

func (this *Sogou) Init(channelUserId, token string, args *map[string]interface{}) {
	this.Lr.Init(channelUserId, token)
	this.args = args
	this.Method = "POST"
	this.Url = "http://dev.app.wan.sogou.com/api/v1/login/verify"
}
	

func (this *Sogou) InitParam() {
	this.Lr.InitParam()

	gid := (*this.args)["SOGOU_GAMEID"].(string)
	appSecret := (*this.args)["SOGOU_APPSECRET"].(string)
	singContext := "gid=%s&session_key=%s&user_id=%s&%s"
	singContext = fmt.Sprintf(singContext,gid,this.token,this.channelUserId,appSecret)
	sign := tool.Md5([]byte(singContext))

	this.Req.Param("gid", gid)
	this.Req.Param("user_id",this.channelUserId)
	this.Req.Param("session_key", this.token)
	this.Req.Param("auth", sign)
}
func (this *Sogou) ParseChannelRet() (err error) {
	// json.Unmarshal([]byte(this.Result), &this.channelRet)
	// if this.channelRet.Error != nil {
	// 	beego.Trace(this.channelRet.Error)
	// }
	beego.Trace(this.Result)
	return json.Unmarshal([]byte(this.Result), &this.channelRet)
}

func (this *Sogou) CheckChannelRet() bool {
	return this.channelRet.Result == true
}
