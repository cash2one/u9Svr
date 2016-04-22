package loginRequest

import (
	"encoding/json"
	// "fmt"
	"u9/models"
	"github.com/astaxie/beego"
	// "u9/tool"
)

//TT语言

type TTChannelRet struct {
	Head    HeadJson `json:"head"`
}

type HeadJson struct{
	Result    string `json:"result"`
	Message   string  `json:"message"`
}

type TT struct {
	Lr
	channelRet TTChannelRet
	channelUserId string 
	token string
	args       *map[string]interface{}
}

 

func LrNewTT(mlr *models.LoginRequest, args *map[string]interface{}) *TT {
	ret := new(TT)
	ret.Init(mlr, args)
	return ret
}

func (this *TT) Init(mlr *models.LoginRequest, args *map[string]interface{}) {
	this.Lr.Init(mlr)
	this.channelUserId = this.mlr.ChannelUserid
	this.token = this.mlr.Token
	this.Method = "POST"
	this.Url = "http://120.132.68.148:18081/sdk.server/rest/user/loginstatus.view"
}

func (this *TT) InitParam() {
	this.Lr.InitParam()
	gameId := (*this.args)["TT_SDK_APPID"].(string)
	// context := `{"uid":"%d","gameId":"%d"}`
	// sign := base64.StdEncoding.EncodeToString(tool.Md5([]byte(context)))

	this.Req.Header("sid", this.token)
	this.Req.Param("uid", this.channelUserId)
	this.Req.Param("gameId", gameId)

}

func (this *TT) ParseChannelRet() (err error) {
	beego.Trace(this.Result)
	return json.Unmarshal([]byte(this.Result), &this.channelRet)
}

func (this *TT) CheckChannelRet() bool {
	return this.channelRet.Head.Result == "0"
}
