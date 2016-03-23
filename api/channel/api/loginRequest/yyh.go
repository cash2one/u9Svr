package loginRequest

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"u9/tool"
)

//当乐

type DangleChannelRet struct {
	Data    DataJson `json:"valid"`
	Status  int      `json:"status"`
	Message string   `json:"message"`
}
type DataJson struct {
	Nick_name    string `json:"nick_name"`
	Valid        string `json:"valid"`
	User_name    string `json:"user_name"`
	Phone        string `json:"phone"`
	Avatar_url   string `json:"avatar_url"`
	Actived      string `json:"actived"`
	Email        string `json:"email"`
	Ticket       string `json:"ticket"`
	Create_tim   string `json:"create_time"`
	User_id      string `json:"user_id"`
	Role_type    string `json:"role_type"`
	Account_type string `json:"account_type"`
}

type Dangle struct {
	Lr
	channelRet DangleChannelRet
}

func LrNewDangle(channelUserId, token string, args *map[string]interface{}) *Dangle {
	ret := new(Dangle)
	ret.Init(channelUserId, token, args)
	return ret
}

func (this *Dangle) Init(channelUserId, token string, args *map[string]interface{}) {
	this.Lr.Init(channelUserId, token)
	loginId := (*args)["YYH_LOGINID"].(string)
	loginKey := (*args)["YYH_LOGINKEY"].(string)
	format := "http://api.appchina.com/appchina-usersdk/user/v2/get.json?"
	+"login_id=%s&login_key=%s&ticket=%s"
	this.Url = fmt.Sprintf(format, loginId, loginKey, token)
}

func (this *Dangle) ParseChannelRet() (err error) {
	beego.Trace(this.Result)
	return json.Unmarshal([]byte(this.Result), &this.channelRet)
}

func (this *Dangle) CheckChannelRet() bool {
	return this.channelRet.Code == 0
}
