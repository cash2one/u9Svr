package loginRequest

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"time"
	"u9/models"
	"u9/tool"
)

//UC九游

type jiuyouReqData struct {
	Sid string `json:"sid"`
}
type jiuyouReqGame struct {
	GameId string `json:"gameId"`
}
type jiuyouChannelReq struct {
	Id   int64         `json:"id"`
	Data jiuyouReqData `json:"data"`
	Game jiuyouReqGame `json:"game"`
	Sign string        `json:"sign"`
}

type jiuyouRetState struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}
type jiuyouRetData struct {
	AccountId string `json:"accountId"`
	Creator   string `json:"creator"`
	NickName  string `json:"nickName"`
}
type jiuyouChannelRet struct {
	Id    uint64         `json:"id"`
	State jiuyouRetState `json:"state"`
	Data  jiuyouRetData  `json:"data"`
}

type JiuYou struct {
	Lr
	channelReq jiuyouChannelReq
	channelRet jiuyouChannelRet
	args       *map[string]interface{}
}

func LrNewJiuYou(mlr *models.LoginRequest, args *map[string]interface{}) *JiuYou {
	ret := new(JiuYou)
	ret.Init(mlr, args)
	return ret
}

func (this *JiuYou) Init(mlr *models.LoginRequest, args *map[string]interface{}) {
	this.Lr.Init(mlr)
	this.Method = "POST"
	this.Url = "http://sdk.g.uc.cn/cp/account.verifySession"
	this.args = args
}

func (this *JiuYou) InitParam() (err error) {
	if err = this.Lr.InitParam(); err != nil {
		return
	}

	appKey := (*this.args)["UC_APPKEY"].(string)

	this.channelReq.Id = time.Now().Unix()
	this.channelReq.Data.Sid = this.mlr.Token
	this.channelReq.Game.GameId = (*this.args)["UC_GAME_ID"].(string)
	this.channelReq.Sign = tool.Md5([]byte("sid=" + this.channelReq.Data.Sid + appKey))

	body, _ := json.Marshal(&this.channelReq)
	beego.Trace(string(body))
	this.Req.Body(string(body))
	return
}

func (this *JiuYou) ParseChannelRet() (err error) {
	beego.Trace(this.Result)
	err = json.Unmarshal([]byte(this.Result), &this.channelRet)
	return
}

func (this *JiuYou) CheckChannelRet() bool {
	return this.channelRet.State.Code == 1
}
