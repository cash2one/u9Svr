package loginRequest

import (
	"fmt"
	"github.com/astaxie/beego"
)

//虫虫
type CCPay struct {
	Lr
}

func LrNewCCPay(channelUserId, token string, args *map[string]interface{}) *CCPay {
	ret := new(CCPay)
	ret.Init(channelUserId, token, args)
	return ret
}

func (this *CCPay) Init(channelUserId, token string, args *map[string]interface{}) {
	this.Lr.Init(channelUserId, token)
	format := "http://android-api.ccplay.com.cn/api/v2/payment/checkUser?token=%s"
	this.Url = fmt.Sprintf(format, token)
}

func (this *CCPay) CheckChannelRet() bool {
	beego.Trace(this.Result)
	return this.Result == "success"
}
