package manager

import (
	"github.com/astaxie/beego"
	"strings"
	"u9/models"
	"u9/www/common"
)

type BaseController struct {
	beego.Controller
	moduleName     string
	controllerName string
	actionName     string
	layoutName     string
}

func (this *BaseController) Prepare() {
	controllerName, actionName := this.GetControllerAndAction()
	this.moduleName = "manager"
	length := len("Controller")
	this.controllerName = strings.ToLower(controllerName[0 : len(controllerName)-length])
	this.actionName = strings.ToLower(actionName)
	this.layoutName = "layout.html"
	this.auth()
}

func (this *BaseController) getAuth() (name, auth string) {
	authData := strings.Split(this.Ctx.GetCookie("managerAuth"), "|")
	if len(authData) == 2 {
		name, auth = authData[0], authData[1]
	}
	return
}

func (this *BaseController) auth() {
	name, auth := this.getAuth()
	if auth == common.GetAuthKey(name) {
		return
	}
	this.Logout()
}

func (this *BaseController) Logout() {
	this.Ctx.SetCookie("managerAuth", "")
	this.Redirect(this.URLFor("LoginController.Login"), 302)
}

func (this *BaseController) display(tpl ...string) {
	var tplname string
	if len(tpl) == 1 {
		tplname = this.moduleName + "/" + tpl[0] + ".html"
	} else {
		tplname = this.moduleName + "/" + this.controllerName + "_" + this.actionName + ".html"
	}

	if this.layoutName != "" {
		this.Layout = this.moduleName + "/" + this.layoutName
	}
	this.TplName = tplname
}

func (this *BaseController) getManager() (manager *models.Manager) {
	manager = new(models.Manager)
	name, _ := this.getAuth()
	manager.Query().Filter("name", name).One(manager)
	return
}

func (this *BaseController) updateData() {
	this.Data["manager"] = this.getManager()
	this.Data["version"] = beego.AppConfig.String("AppVer")
}

func (this *BaseController) showMsg(msg ...string) {
	if len(msg) == 1 {
		msg = append(msg, this.Ctx.Request.Referer())
	}

	this.Data["msg"] = msg[0]
	this.Data["redirect"] = msg[1]
	this.Layout = this.moduleName + "/layout.html"
	this.TplName = "common" + "/" + "showmsg.html"
	this.Render()
	this.StopRun()
}

func (this *BaseController) getClientIp() string {
	s := strings.Split(this.Ctx.Request.RemoteAddr, ":")
	return s[0]
}
