package manager

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"strings"
	"u9/models"
	"u9/www/common"
)

type CpController struct {
	BaseController
}

//用户列表
func (this *CpController) List() {
	var page int64
	var pagesize int64 = 10
	var list []*models.Cp
	var cp models.Cp

	if page, _ = this.GetInt64("page"); page < 1 {
		page = 1
	}
	offset := (page - 1) * pagesize

	count, _ := cp.Query().Count()
	if count > 0 {
		cp.Query().OrderBy("-id").Limit(pagesize, offset).All(&list)
	}

	this.Data["list"] = list
	pageBar := common.NewPager(page, count, pagesize, "/manager/cp/list?page=%d").ToString()
	this.Data["pagebar"] = pageBar

	this.updateData()
	this.display()
}

//添加用户
func (this *CpController) Add() {
	input := make(map[string]string)
	errMsg := make(map[string]string)
	if this.Ctx.Request.Method == "POST" {
		name := strings.TrimSpace(this.GetString("Name"))
		password1 := strings.TrimSpace(this.GetString("Password1"))
		password2 := strings.TrimSpace(this.GetString("Password2"))
		email := strings.TrimSpace(this.GetString("Email"))
		state, _ := this.GetInt64("State")

		input["Name"] = name
		input["Password1"] = password1
		input["Password2"] = password2
		input["Email"] = email

		valid := validation.Validation{}

		if v := valid.Required(name, "Name"); !v.Ok {
			errMsg["Name"] = "请输入用户名"
		} else if v := valid.MaxSize(name, 15, "Name"); !v.Ok {
			errMsg["Name"] = "用户名长度不能大于15个字符"
		}

		if v := valid.Required(password1, "password"); !v.Ok {
			errMsg["Password1"] = "请输入密码"
		}

		if v := valid.Required(password2, "password2"); !v.Ok {
			errMsg["Password2"] = "请再次输入密码"
		} else if password1 != password2 {
			errMsg["Password2"] = "两次输入的密码不一致"
		}

		if v := valid.Required(email, "email"); !v.Ok {
			errMsg["Email"] = "请输入email地址"
		} else if v := valid.Email(email, "email"); !v.Ok {
			errMsg["Email"] = "Email无效"
		}

		if state > 0 {
			state = 1
		} else {
			state = 0
		}

		if len(errMsg) == 0 {
			var cp models.Cp
			cp.Name = name
			cp.Password = password1 //models.Md5([]byte(password1))
			cp.Email = email
			cp.State = int8(state)
			if err := cp.Insert(); err != nil {
				this.showMsg(err.Error())
			}
			this.Redirect("/manager/cp/list", 302)
		}

	}

	this.Data["input"] = input
	this.Data["errMsg"] = errMsg
	this.updateData()
	this.display()
}

//编辑用户
func (this *CpController) Edit() {
	id, _ := this.GetInt("id")
	cp := models.Cp{Id: id}
	if err := cp.Read(); err != nil {
		this.showMsg("用户不存在")
	}

	errMsg := make(map[string]string)

	if this.Ctx.Request.Method == "POST" {
		password1 := strings.TrimSpace(this.GetString("password1"))
		password2 := strings.TrimSpace(this.GetString("password2"))
		email := strings.TrimSpace(this.GetString("email"))
		state, _ := this.GetInt64("state")
		valid := validation.Validation{}

		if v := valid.Required(password1, "password1"); !v.Ok {
			errMsg["password1"] = "请输入密码"
			beego.Trace(errMsg)
		} else if password1 != password2 {
			errMsg["password2"] = "两次输入密码不一致"
			beego.Trace(errMsg)
		} else {
			cp.Password = password1 //models.Md5([]byte(password))
		}

		if v := valid.Required(email, "email"); !v.Ok {
			errMsg["email"] = "请输入email地址"
		} else if v := valid.Email(email, "email"); !v.Ok {
			errMsg["email"] = "Email无效"
		} else {
			cp.Email = email
		}

		if state > 0 {
			cp.State = 1
		} else {
			cp.State = 0
		}

		if len(errMsg) == 0 {
			cp.Update()
			this.Redirect("/manager/cp/list", 302)
		}
	}
	this.Data["errMsg"] = errMsg
	this.Data["cp"] = cp
	this.updateData()
	this.display()
}

//删除用户
func (this *CpController) Delete() {
	id, _ := this.GetInt("id")
	cp := models.Cp{Id: id}
	if cp.Read() == nil {
		cp.Delete()
	}

	this.Redirect("/manager/cp/list", 302)
}
