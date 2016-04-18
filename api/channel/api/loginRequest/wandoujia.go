package loginRequest

import (
	"fmt"
	"u9/models"
)

type Wandoujia struct {
	Lr
}

func LrNewWandoujia(mlr *models.LoginRequest, args *map[string]interface{}) *Wandoujia {
	ret := new(Wandoujia)
	ret.Init(mlr, args)
	return ret
}

func (this *Wandoujia) Init(mlr *models.LoginRequest, args *map[string]interface{}) {
	this.Lr.Init(mlr)
	this.IsHttps = true
	format := "https://pay.wandoujia.com/api/uid/check?uid=%s&token=%s&appkey_id=%s"
	appkey := (*args)["WANDOUJIA_APPKEY"].(string)
	this.Url = fmt.Sprintf(format, this.mlr.ChannelUserid, this.mlr.Token, appkey)
}

func (this *Wandoujia) CheckChannelRet() bool {
	return this.Result != "false"
}
