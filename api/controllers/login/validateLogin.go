package login

import (
	"fmt"
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

func (this *LoginController) ValidateLogin() {
	ret := new(common.BasicRet).Init()

	defer func() {
		this.Data["json"] = ret
		this.ServeJSON(true)
	}()

	msg := common.DumpCtx(this.Ctx)
	beego.Trace(msg)

	beego.Trace(`validateLogin: 1:validate`)
	vlp := new(ValidateLoginParam)
	if code := this.Validate(vlp); code != 0 {
		ret.SetCode(code)
		return
	}

	beego.Trace("validateLogin: 2:callLoginRequest")
	ret = channelApi.CallLoginRequest(vlp.lr)
	if ret.Code == 0 {
		beego.Trace("validateLogin: 3:updateLoginLog")
		vlp.updateLoginLog()
	} else {
		beego.Error("validateLogin: ret.Code!=0")
	}
}

type ValidateLoginParam struct {
	UserId string `json:"ChannelUserId"`
	Token  string `json:"Token"`
	lr     *models.LoginRequest
}

func (this *ValidateLoginParam) handleSpecialToken() {
	if this.lr.ChannelId == 134 && this.lr.ProductId == 1002 {
		beego.Trace(`handleSpecialToken: 暴走水浒 and 芒果玩`)
		this.Token = strings.Replace(this.Token, `"`, ``, -1)
	}
}

func (this *ValidateLoginParam) Valid(v *validation.Validation) {
	switch {
	case strings.TrimSpace(this.UserId) == "":

		msg := `valid: require userId`

		beego.Error(msg)
		v.SetError("1004", msg)
		return
	case strings.TrimSpace(this.Token) == "":

		msg := `valid: require token`
		beego.Error(msg)

		v.SetError("1003", msg)
		return
	}

	beego.Trace(`valid: 1:check userId`)
	this.lr = new(models.LoginRequest)
	qs := this.lr.Query().Filter("userId", this.UserId)
	if err := qs.One(this.lr); err != nil {

		format := `valid: err:%v`
		msg := fmt.Sprintf(format, err)
		beego.Error(msg)

		v.SetError("1004", msg)
		return
	}

	beego.Trace(`valid: 2:handleSpecialToken`)
	this.handleSpecialToken()

	beego.Trace(`valid: 3:check token`)
	if qs.Filter("token", this.Token).Exist() == false {

		this.Token, _ = url.QueryUnescape(this.Token)

		if qs.Filter("token", this.Token).Exist() == false {
			msg := `valid: token isn't exist`
			beego.Error(msg)
			v.SetError("1003", msg)
		}
		return
	}
}

func (this *ValidateLoginParam) updateLoginLog() {
	lrl := models.LoginRequestLog{
		LoginRequestId: this.lr.Id,
		LoginTime:      time.Now()}
	if _, _, err := orm.NewOrm().ReadOrCreate(&lrl,
		"LoginRequestId", "LoginTime"); err != nil {
		beego.Warn(err)
	}
}
