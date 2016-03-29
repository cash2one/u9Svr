package loginRequest

import (
	"fmt"
	// "github.com/astaxie/beego"
	"strconv"
	"time"
	"u9/tool"
)

//叉叉

type Guopan struct {
	Lr
}

func LrNewGuopan(channelUserId, token string, args *map[string]interface{}) *Guopan {
	ret := new(Guopan)
	ret.Init(channelUserId, token, args)
	return ret
}

func (this *Guopan) Init(channelUserId, token string, args *map[string]interface{}) {
	this.Lr.Init(channelUserId, token)
	appid := (*args)["GUOPAN_APPID"].(string)
	secretKey := (*args)["GUOPAN_SERVER_SECRETKEY"].(string)
	t := strconv.FormatInt(time.Now().Unix(), 10)
	context := channelUserId + appid + t + secretKey
	sign := tool.Md5([]byte(context))
	format := "http://userapi.guopan.cn/gamesdk/verify?game_uin=%s&appid=%s&token=%s&t=%s&sign=%s"
	this.Url = fmt.Sprintf(format, channelUserId, appid, token, t, sign)
}

func (this *Guopan) CheckChannelRet() bool {
	return this.Result == "true"
}
