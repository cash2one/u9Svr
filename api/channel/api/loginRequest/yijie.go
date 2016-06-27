package loginRequest

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"strings"
	"u9/models"
)

//易接

type YiJie struct {
	Lr
}

func LrNewYiJie(mlr *models.LoginRequest, args *map[string]interface{}) *YiJie {
	ret := new(YiJie)
	ret.Init(mlr, args)
	return ret
}

func (this *YiJie) Init(mlr *models.LoginRequest, args *map[string]interface{}) {
	this.Lr.Init(mlr)

	publishChannelId := ""
	var clientParam map[string]string
	var err error
	if err = json.Unmarshal([]byte(mlr.Ext), &clientParam); err != nil {
		ok := false
		publishChannelId, ok = (*args)["com.snowfish.channelid"].(string)
		if !ok {
			beego.Error("init: publishChannelId isn't exist.")
			return
		}
	} else {
		publishChannelId = clientParam["channelId"]
	}

	publishChannelId = strings.Replace(publishChannelId, "{", "", -1)
	publishChannelId = strings.Replace(publishChannelId, "}", "", -1)
	publishChannelId = strings.Replace(publishChannelId, "-", "", -1)

	appid := (*args)["com.snowfish.appid"].(string)
	appid = strings.Replace(appid, "{", "", -1)
	appid = strings.Replace(appid, "}", "", -1)
	appid = strings.Replace(appid, "-", "", -1)

	format := "http://sync.1sdk.cn/login/check.html?sdk=%s&app=%s&uin=%s&sess=%s"
	this.Url = fmt.Sprintf(format,
		publishChannelId,
		appid,
		this.mlr.ChannelUserid,
		this.mlr.Token)
}

func (this *YiJie) CheckChannelRet() (ret bool) {
	ret = this.Result == "0"
	format := "checkChannelRet: result: %s, url:%s"
	msg := fmt.Sprintf(format, this.Result, this.Url)

	if !ret {
		beego.Error(msg)
	} else {
		beego.Trace(msg)
	}
	return
}
