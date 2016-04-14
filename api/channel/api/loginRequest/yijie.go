package loginRequest

import (
	"fmt"
	"github.com/astaxie/beego"
	"strings"
)

//易接

type YiJie struct {
	Lr
	channelRet int
}

func LrNewYiJie(channelUserId, token string, args *map[string]interface{}) *YiJie {
	ret := new(YiJie)
	ret.Init(channelUserId, token, args)
	return ret
}

func (this *YiJie) Init(channelUserId, token string, args *map[string]interface{}) {
	this.Lr.Init(channelUserId, token)
	appid := replaceSDKParam((*args)["com.snowfish.appid"].(string))
	channleid := replaceSDKParam((*args)["com.snowfish.channelid"].(string))
	format := "http://sync.1sdk.cn/login/check.html?sdk=%s&app=%s&uin=%s&sess=%s"
	this.Url = fmt.Sprintf(format, channleid, appid, channelUserId, token)
}

func (this *YiJie) ParseChannelRet() (err error) {
	beego.Trace(this.Result)
	return
}

func (this *YiJie) CheckChannelRet() bool {
	return this.channelRet == 0
}

func replaceSDKParam(s string) (ret string) {
	ret = strings.Replace(s, "{", "", -1)
	ret = strings.Replace(ret, "}", "", -1)
	ret = strings.Replace(ret, "-", "", -1)
	return
}