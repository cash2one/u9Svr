package manager

import (
	//"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"strings"
	"u9/models"
	"u9/www/common"
)

type ChannelController struct {
	BaseController
}

func (this *ChannelController) List() {
	var page int64
	var pagesize int64 = 10
	var list []*models.Channel
	var channel models.Channel

	if page, _ = this.GetInt64("page"); page < 1 {
		page = 1
	}
	offset := (page - 1) * pagesize

	count, _ := channel.Query().Count()
	if count > 0 {
		channel.Query().OrderBy("-id").Limit(pagesize, offset).All(&list)
	}

	this.Data["list"] = list
	pageBar := common.NewPager(page, count, pagesize, "/manager/channel/list?page=%d").ToString()
	this.Data["pagebar"] = pageBar

	this.updateData()
	this.display()
}

func (this *ChannelController) Add() {
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
			var channel models.Channel
			channel.Name = name
			if err := channel.Insert(); err != nil {
				this.showMsg(err.Error())
			}
			this.Redirect("/manager/channel/list", 302)
		}

	}

	this.Data["input"] = input
	this.Data["errMsg"] = errMsg
	this.updateData()
	this.display()
}

func (this *ChannelController) Edit() {
	id, _ := this.GetInt("id")
	channel := models.Channel{Id: id}
	if err := channel.Read(); err != nil {
		this.showMsg("用户不存在")
	}

	errMsg := make(map[string]string)

	if this.Ctx.Request.Method == "POST" {
		password := strings.TrimSpace(this.GetString("password"))
		password2 := strings.TrimSpace(this.GetString("password2"))
		email := strings.TrimSpace(this.GetString("email"))
		state, _ := this.GetInt64("state")
		valid := validation.Validation{}

		if password != "" {
			if v := valid.Required(password2, "password2"); !v.Ok {
				errMsg["password2"] = "请再次输入密码"
			} else if password != password2 {
				errMsg["password2"] = "两次输入的密码不一致"
			} else {
				//channel.Password = models.Md5([]byte(password))
			}
		}
		if v := valid.Required(email, "email"); !v.Ok {
			errMsg["email"] = "请输入email地址"
		} else if v := valid.Email(email, "email"); !v.Ok {
			errMsg["email"] = "Email无效"
		} else {
			//channel.Email = email
		}

		if state > 0 {
			//channel.State = 1
		} else {
			//channel.State = 0
		}

		if len(errMsg) == 0 {
			channel.Update()
			this.Redirect("/manager/channel/list", 302)
		}
	}
	this.Data["errMsg"] = errMsg
	this.Data["channel"] = channel
	this.updateData()
	this.display()
}

func (this *ChannelController) Delete() {
	id, _ := this.GetInt("id")
	channel := models.Channel{Id: id}
	if channel.Read() == nil {
		channel.Delete()
	}

	this.Redirect("/manager/channel/list", 302)
}
