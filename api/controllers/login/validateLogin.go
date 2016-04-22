package login

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	"net/url"
	"strings"
	"time"
	"u9/api/channel/api"
	"u9/api/common"
	"u9/models"
)

type ValidateLoginParam struct {
	UserId string `json:"ChannelUserId"`
	Token  string `json:"Token"`
	lr     *models.LoginRequest
	lrl    *models.LoginRequestLog
}

func (this *ValidateLoginParam) handleSpecialToken() {
	//暴走水浒 and 芒果玩
	if this.lr.ChannelId == 134 && this.lr.ProductId == 1002 {
		this.Token = strings.Replace(this.Token, `"`, ``, -1)
	}
}

func (this *ValidateLoginParam) Valid(v *validation.Validation) {
	switch {
	case strings.TrimSpace(this.UserId) == "":
		v.SetError("1004", "Require userId")
		return
	case strings.TrimSpace(this.Token) == "":
		v.SetError("1003", "Require token")
		return
	}

	this.lr = new(models.LoginRequest)
	qs := this.lr.Query().Filter("userId", this.UserId)
	if err := qs.One(this.lr); err != nil {
		//beego.Error(this.Input())
		v.SetError("1004", "Record isn't exist in table:loginRequest with UserId="+this.UserId)
		return
	}

	this.handleSpecialToken()
	if qs.Filter("token", this.Token).Exist() == false {
		this.Token, _ = url.QueryUnescape(this.Token)
		if qs.Filter("token", this.Token).Exist() == false {
			//beego.Error(this.Input())
			v.SetError("1003", "Record isn't exist in table:loginRequest with token:"+this.Token)
		}
		return
	}

}

func (this *LoginController) ValidateLogin() {
	ret := new(common.BasicRet).Init()

	defer func() {
		this.Data["json"] = ret
		this.ServeJSON(true)
	}()

	vlp := new(ValidateLoginParam)
	if code := this.Validate(vlp); code != 0 {
		ret.SetCode(code)
		return
	}
	ret = channelApi.CallLoginRequest(vlp.lr)
	if ret.Code == 0 {
		vlp.addDB()
	}

}

func (this *ValidateLoginParam) addDB() {
	lr := models.LoginRequestLog{
		LoginRequestId: this.lr.Id,
		LoginTime:      time.Now()}
	if _, _, err := orm.NewOrm().ReadOrCreate(&lr, "LoginRequestId", "LoginTime"); err != nil {
		beego.Trace(err)
	}
}
