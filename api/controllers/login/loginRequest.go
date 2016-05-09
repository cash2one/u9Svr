package login

import (
	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"

	"time"
	"u9/api/common"
	"u9/api/controllers/login/loginRequestHandle"
	"u9/models"
)

type LoginRequestRet struct {
	common.BasicRet
	UserId     string `json:"UserId"`
	TransType  string `json:"TransType"`
	ChannelExt string `json:"Ext"`
}

func (this *LoginRequestRet) Init() *LoginRequestRet {
	this.BasicRet.Init()
	this.TransType = "CREATE_USER"
	return this
}

func (this *LoginController) handleLrParam() (err error) {
	if code := this.Validate(&this.lrParam); code != 0 {
		this.lrRet.SetCode(code)
		err = errors.New("loginRequest' param parse error.")
		return
	}
	switch this.lrParam.ChannelId {
	case 123: //熊猫玩
		var lrh loginRequestHandle.LRHandle
		lrh = loginRequestHandle.NewXMW()
		if err = lrh.Init(&this.lrParam); err != nil {
			return
		}
		if this.lrRet.ChannelExt, err = lrh.Handle(); err != nil {
			return
		}
		this.lrParam.Token = lrh.GetToken()
		return
	default:
		if this.lrParam.ChannelUserId == "" {
			this.lrRet.SetCode(1001)
			err = errors.New("Require channelUserId.")
			return
		}
	}
	return
}

func (this *LoginController) updateDB() (err error) {
	userId := models.GenerateUserId(
		this.lrParam.ChannelId,
		this.lrParam.ProductId,
		this.lrParam.ChannelUserId)

	lr := models.LoginRequest{
		ChannelId:       this.lrParam.ChannelId,
		ProductId:       this.lrParam.ProductId,
		ChannelUserid:   this.lrParam.ChannelUserId,
		Token:           this.lrParam.Token,
		IsDebug:         this.lrParam.IsDebug,
		ChannelUsername: this.lrParam.ChannelUserName,
		Ext:             this.lrParam.Ext,
		Userid:          userId,
		UpdateTime:      time.Now()}

	create := false
	if create, _, err = orm.NewOrm().ReadOrCreate(&lr,
		"ChannelId", "ProductId", "ChannelUserid", "Userid"); err != nil {
		beego.Error(lr)
		return
	}

	if !create {
		lr.Token = this.lrParam.Token
		lr.ChannelUsername = this.lrParam.ChannelUserName
		lr.Ext = this.lrParam.Ext
		//beego.Trace("loginUrlExt:", lr.Ext)
		lr.UpdateTime = time.Now()
		if err = lr.Update("ChannelUsername", "Token", "IsDebug", "UpdateTime", "Ext"); err != nil {
			beego.Error(lr)
			return
		}
	}
	this.lrRet.UserId = lr.Userid
	return
}

func (this *LoginController) LoginRequest() {
	this.lrRet.Init()

	defer func() {
		this.Data["json"] = this.lrRet
		this.ServeJSON(true)
	}()

	if err := this.handleLrParam(); err != nil {
		beego.Error(err)
		return
	}

	if err := this.updateDB(); err != nil {
		beego.Error(err)
		return
	}

	this.lrRet.SetCode(0)
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
