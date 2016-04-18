package loginRequest

import (
	"fmt"
	"github.com/astaxie/beego"
	"u9/models"
)

//虫虫
type CCPay struct {
	Lr
}

func LrNewCCPay(mlr *models.LoginRequest, args *map[string]interface{}) *CCPay {
	ret := new(CCPay)
	ret.Init(mlr, args)
	return ret
}

func (this *CCPay) Init(mlr *models.LoginRequest, args *map[string]interface{}) {
	this.Lr.Init(mlr)
	format := "http://android-api.ccplay.com.cn/api/v2/payment/checkUser?token=%s"
	this.Url = fmt.Sprintf(format, this.mlr.Token)
}

func (this *CCPay) CheckChannelRet() bool {
	beego.Trace(this.Result)
	return this.Result == "success"
}
