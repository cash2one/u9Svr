package loginRequest

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"u9/models"
	"u9/tool"
)

//微游汇

type WeiUUChannelRet struct {
	State    int 		`json:"state"`
	Data     dataJson   `json:"data"`
}
type dataJson struct{
	UserID 		int 	`json:userID`
	Username 	string  `json:username`
}

type WeiUU struct {
	Lr
	channelRet WeiUUChannelRet
}

func LrNewWeiUU(mlr *models.LoginRequest, args *map[string]interface{}) *WeiUU {
	ret := new(WeiUU)
	ret.Init(mlr, args)
	return ret
}

func (this *WeiUU) Init(mlr *models.LoginRequest, args *map[string]interface{}) {
	this.Lr.Init(mlr)
	// appid := (*args)["DANGLE_SDK_APPID"].(string)
	appkey := (*args)["WEIUU_APPKEY"].(string)
	signFormat := "userID=%stoken=%s%s"
	signContent := fmt.Sprintf(signFormat, mlr.ChannelUserid , mlr.Token ,appkey)
	sign := tool.Md5([]byte(signContent))
	format := "http://unionsdk.weiuu.cn/user/loginServer?userID=%s&token=%s&sign=%s"
	this.Url = fmt.Sprintf(format, mlr.ChannelUserid ,mlr.Token, sign)
 	beego.Trace(this.Url)
}

func (this *WeiUU) ParseChannelRet() (err error) {
	beego.Trace(this.Result)
	return json.Unmarshal([]byte(this.Result), &this.channelRet)
}

func (this *WeiUU) CheckChannelRet() bool {
	return this.channelRet.State == 1
}
