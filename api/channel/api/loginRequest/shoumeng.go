package loginRequest

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"u9/models"
)

//手盟

type ShouMengChannelRet struct {
	Result  int    `json:"result"`
	Message string `json:"message"`
}

type TokenJson struct {
	Login_account string `json:"login_account"`
	Session_id    string `json:"session_id"`
}

type ShouMeng struct {
	Lr
	tokenJson  TokenJson
	channelRet ShouMengChannelRet
}

func LrNewShouMeng(mlr *models.LoginRequest, args *map[string]interface{}) *ShouMeng {
	ret := new(ShouMeng)
	ret.Init(mlr, args)
	return ret
}

func (this *ShouMeng) Init(mlr *models.LoginRequest, args *map[string]interface{}) {
	this.Lr.Init(mlr)
	this.Method = "POST"
	// appid := (*args)["DANGLE_SDK_APPID"].(string)
	// appkey := (*args)["DANGLE_SDK_APPKEY"].(string)
	json.Unmarshal([]byte(this.mlr.Token), &this.tokenJson)
	this.Url = "http://www.19meng.com/api/v1/verify_session_id"
	beego.Trace(this.Url)
}
func (this *ShouMeng) InitParam() {
	this.Lr.InitParam()
	this.Req.Param("user_id", this.mlr.ChannelUserid)
	this.Req.Param("login_account", this.tokenJson.Login_account)
	this.Req.Param("session_id", this.tokenJson.Session_id)
}
func (this *ShouMeng) ParseChannelRet() (err error) {
	beego.Trace(this.Result)
	return json.Unmarshal([]byte(this.Result), &this.channelRet)
}

func (this *ShouMeng) CheckChannelRet() bool {
	return this.channelRet.Result == 1
}
