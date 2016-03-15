package manager

import (
	"github.com/astaxie/beego/validation"
	"u9/models"
)

type ProfileParam struct {
	manager      *models.Manager
	OldPassword  string
	NewPassword1 string
	NewPassword2 string
}

func (this *ProfileParam) Valid(v *validation.Validation) {
	switch {
	case len(this.NewPassword1) < 6:
		v.SetError("newPassword1", "密码长度不能少于6个字符.")
		return
	case this.NewPassword1 != this.NewPassword2:
		v.SetError("newPassword2", "两次输入的密码不一致.")
		return
	}

	if this.manager.Password != this.OldPassword {
		v.SetError("password", "当前密码错误.")
		return
	}
}

type ProfileController struct {
	BaseController
	profileParam ProfileParam
	errMsg       map[string]string
}

func (this *ProfileController) checkParam() bool {
	this.profileParam.manager = this.getManager()
	if err := this.ParseForm(&this.profileParam); err != nil {
		this.showMsg(err.Error())
	}

	valid := validation.Validation{}
	if valid.Valid(&(this.profileParam)); valid.HasErrors() {
		this.errMsg = make(map[string]string)
		for _, err := range valid.Errors {
			this.errMsg[err.Key] = err.Message
			return false
		}
	}
	return true
}

func (this *ProfileController) updateData() {
	manager := this.getManager()
	manager.Password = this.profileParam.NewPassword1
	manager.Update("password")
	this.Data["updated"] = true
}

func (this *ProfileController) Profile() {
	if this.Ctx.Request.Method == "POST" {
		if this.checkParam() {
			this.updateData()
		}
	}
	this.BaseController.updateData()
	this.Data["errMsg"] = this.errMsg
	this.display()
}
