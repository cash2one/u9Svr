package loginRequest

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	// "u9/tool"
	"u9/models"
)

//pptv

type PPTVChannelRet struct {
	Status  int    `json:"status"`
	Message string `json:"msg"`
}

type PPTV struct {
	Lr
	channelRet PPTVChannelRet
}

func LrNewPPTV(mlr *models.LoginRequest, args *map[string]interface{}) *PPTV {
	ret := new(PPTV)
	ret.Init(mlr, args)
	return ret
}

func (this *PPTV) Init(mlr *models.LoginRequest, args *map[string]interface{}) {
	this.Lr.Init(mlr)
	// appid := (*args)["PPTV_SDK_APPID"].(string)
	// appkey := (*args)["PPTV_SDK_APPKEY"].(string)
	// sign := tool.Md5([]byte(appid + "|" + appkey + "|" + token + "|" + this.mlr.ChannelUserid))
	format := "http://api.user.vas.pptv.com/c/v2/cksession.php?type=login&sessionid=%s&username=%s&app=mobgame"
	this.Url = fmt.Sprintf(format, this.mlr.Token, this.mlr.ChannelUserid)
}

func (this *PPTV) ParseChannelRet() (err error) {
	beego.Trace(this.Result)
	return json.Unmarshal([]byte(this.Result), &this.channelRet)
}

func (this *PPTV) CheckChannelRet() bool {
	return this.channelRet.Status == 1
}
