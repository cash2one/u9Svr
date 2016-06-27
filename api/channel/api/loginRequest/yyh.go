package loginRequest

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"u9/models"
)

//应用汇
type yyhChannelRet struct {
	Data struct {
		Nick_name    string `json:"nick_name"`
		Valid        bool   `json:"valid"`
		User_name    string `json:"user_name"`
		Phone        string `json:"phone"`
		Avatar_url   string `json:"avatar_url"`
		Actived      bool   `json:"actived"`
		Email        string `json:"email"`
		Ticket       string `json:"ticket"`
		Create_time  string `json:"create_time"`
		User_id      int    `json:"user_id"`
		Role_type    int    `json:"role_type"`
		Account_type string `json:"account_type"`
	} `json:"data"`
	Status  int    `json:"status`
	Message string `json:"message`
}

type YYH struct {
	Lr
	channelRet yyhChannelRet
}

func LrNewYYH(mlr *models.LoginRequest, args *map[string]interface{}) *YYH {
	ret := new(YYH)
	ret.Init(mlr, args)
	return ret
}

func (this *YYH) Init(mlr *models.LoginRequest, args *map[string]interface{}) {
	this.Lr.Init(mlr)
	loginId := (*args)["YYH_LOGINID"].(string)
	loginKey := (*args)["YYH_LOGINKEY"].(string)
	format := "http://api.appchina.com/appchina-usersdk/user/v2/get.json?login_id=%s&login_key=%s&ticket=%s"
	this.Url = fmt.Sprintf(format, loginId, loginKey, this.mlr.Token)
}

func (this *YYH) ParseChannelRet() (err error) {
	beego.Trace(this.Result)
	return json.Unmarshal([]byte(this.Result), &this.channelRet)
}

func (this *YYH) CheckChannelRet() bool {
	return this.channelRet.Status == 0
}
