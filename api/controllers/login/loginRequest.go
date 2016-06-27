package login

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
	"u9/api/common"
	"u9/api/controllers/login/lrBefore"
	"u9/models"
)

type LoginRequestRet struct {
	common.BasicRet
	UserId     string    `json:"UserId"`
	TransType  string    `json:"TransType"`
	CreateTime time.Time `json:"CreateTime"`
	ChannelExt string    `json:"Ext"`
}

func (this *LoginRequestRet) Init() *LoginRequestRet {
	this.BasicRet.Init()
	this.TransType = "CREATE_USER"
	return this
}

func (this *LoginController) beforeHandle() (err error) {
	if code := this.Validate(&this.lrParam); code != 0 {
		this.lrRet.SetCode(code)
		format := "beforeHandle: param is error"
		msg := fmt.Sprintf(format)
		err = errors.New(msg)
		beego.Error(err)
		return
	}

	var bh lrBefore.Handle
	switch this.lrParam.ChannelId {
	case 102: //qihoo360
		bh = lrBefore.NewQihoo360()
	case 123: //熊猫玩
		bh = lrBefore.NewXMW()
	case 145: //huawei
		bh = lrBefore.NewHuawei()
	case 146: //lenovo
		bh = lrBefore.NewLenovo()
	case 149: //coolPad
		bh = lrBefore.NewCoolPad()
	default:
		if this.lrParam.ChannelUserId == "" {
			this.lrRet.SetCode(1001)
			format := `beforeHandle: default require channelUserId`
			msg := fmt.Sprintf(format)
			err = errors.New(msg)
			beego.Error(err)
		}
		return
	}

	beego.Trace("beforeHandle:1:Init")
	if err = bh.Init(&this.lrParam); err != nil {
		return
	}

	beego.Trace("beforeHandle:2:Exec")
	if this.lrRet.ChannelExt, err = bh.Exec(); err != nil {
		return
	}
	return
}

func (this *LoginController) updateDB() (err error) {
	userId := models.GenerateUserId(
		this.lrParam.ChannelId,
		this.lrParam.ProductId,
		this.lrParam.ChannelUserId)

	this.lr = models.LoginRequest{
		ChannelId:     this.lrParam.ChannelId,
		ProductId:     this.lrParam.ProductId,
		ChannelUserid: this.lrParam.ChannelUserId,
		Userid:        userId,
		MobileInfo:    this.lrParam.MobileInfo,

		Token:           this.lrParam.Token,
		ChannelUsername: this.lrParam.ChannelUserName,
		Ext:             this.lrParam.Ext,
		IsDebug:         this.lrParam.IsDebug,
		UpdateTime:      time.Now()}

	create := false
	if create, _, err = orm.NewOrm().ReadOrCreate(&this.lr,
		"ChannelId", "ProductId", "ChannelUserid", "Userid"); err != nil {
		format := `updateDB: err:%v`
		msg := fmt.Sprintf(format, err)
		beego.Error(msg)
		return
	}

	if !create {
		this.lr.Token = this.lrParam.Token
		this.lr.ChannelUsername = this.lrParam.ChannelUserName
		this.lr.Ext = this.lrParam.Ext
		this.lr.IsDebug = this.lrParam.IsDebug
		this.lr.UpdateTime = time.Now()

		if err = this.lr.Update("ChannelUsername", "Token", "IsDebug",
			"UpdateTime", "MobileInfo", "Ext"); err != nil {
			format := `updateDB: err:%v`
			msg := fmt.Sprintf(format, err)
			beego.Error(msg)
			return
		}
	}
	this.lrRet.UserId = this.lr.Userid
	return
}

func (this *LoginController) LoginRequest() {
	this.lrRet.Init()

	msg := "loginRequest:" + common.DumpCtx(this.Ctx)
	beego.Trace(msg)

	defer func() {
		this.Data["json"] = this.lrRet
		this.ServeJSON(true)
	}()

	beego.Trace("loginRequest: 1:beforeHandle")
	if err := this.beforeHandle(); err != nil {
		return
	}

	beego.Trace("loginRequest: 2:updateDB")
	if err := this.updateDB(); err != nil {
		beego.Error(err)
		return
	}

	this.lrRet.SetCode(0)

	this.lrRet.CreateTime = this.lr.CreateTime
}

/*
  游戏登录请求
  test url:
  http://192.168.0.185/api/gameLoginRequest/?
  ProductId=1000&
  UserId=23c16f5323755132272fba79ab2e11d8&
  ProductOrderId=game20160114142841787&
  Amount=32&
  CallbackUrl=http://www.baidu.com
*/
