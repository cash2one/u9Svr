package loginRequest

import (
	"fmt"
)

type Wandoujia struct {
	Lr
}

func LrNewWandoujia(channelUserId, token string, args *map[string]interface{}) *Wandoujia {
	ret := new(Wandoujia)
	ret.Init(channelUserId, token, args)
	return ret
}

func (this *Wandoujia) Init(channelUserId, token string, args *map[string]interface{}) {
	this.Lr.Init(channelUserId, token)
	this.IsHttps = true
	format := "https://pay.wandoujia.com/api/uid/check?uid=%s&token=%s&appkey_id=%s"
	appkey := (*args)["WANDOUJIA_APPKEY"].(string)
	this.Url = fmt.Sprintf(format, channelUserId, token, appkey)
}

func (this *Wandoujia) CheckChannelRet() bool {
	return this.Result != "false"
}
