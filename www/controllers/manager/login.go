package manager

import (
	"github.com/astaxie/beego/validation"
	"strings"
	"time"
	"u9/models"
	"u9/www/common"
)

type LoginParam struct {
	Name     string
	Password string
	Remember string
}

func (this *LoginParam) Valid(v *validation.Validation) {
	switch {
	case strings.TrimSpace(this.Name) == "":
		v.SetError("login", "需要登录名称.")
		return
	case strings.TrimSpace(this.Password) == "":
		v.SetError("login", "需要登录密码.")
		return
	}

	isExist := new(models.Manager).Query().
		Filter("name", this.Name).Filter("password", this.Password).Exist()
	if !isExist {
		v.SetError("login", "不存在该账号或者密码错误.")
	}
}

type LoginController struct {
	BaseController
	loginParam LoginParam
}

func (this *LoginController) Prepare() {

}

func (this *LoginController) checkParam() bool {
	if err := this.ParseForm(&this.loginParam); err != nil {
		this.Data["errMsg"] = err.Error()
		return false
	}
	valid := validation.Validation{}
	if valid.Valid(&(this.loginParam)); valid.HasErrors() {
		for _, vErr := range valid.Errors {
			this.Data["errMsg"] = vErr.Message
			return false
		}
	}
	return true
}

func (this *LoginController) updateData() {
	manager := new(models.Manager)
	manager.Query().Filter("name", this.loginParam.Name).One(manager)

	manager.LoginCount += 1
	manager.LastLoginIp = this.getClientIp()
	manager.LastLoginTime = time.Now()
	manager.Update()

	auth := manager.Name + "|" + common.GetAuthKey(manager.Name)
	if this.loginParam.Remember == "yes" {
		this.Ctx.SetCookie("managerAuth", auth, 7*86400)
	} else {
		this.Ctx.SetCookie("managerAuth", auth)
	}
}

func (this *LoginController) Login() {
	if this.Ctx.Request.Method == "POST" {
		if this.checkParam() {
			this.updateData()
			this.Redirect(this.URLFor("IndexController.Index"), 302)
			return
		}
	}
	this.TplName = "manager/login.html"
}
