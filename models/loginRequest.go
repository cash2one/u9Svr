package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"time"
	"u9/tool"
)

type LoginRequest struct {
	Id              int
	ChannelId       int
	ChannelUsername string
	ChannelUserid   string
	Token           string
	ProductId       int
	IsDebug         bool
	Userid          string
	UpdateTime      time.Time `orm:"auto_now;type(datatime)"`
	//Channel         *Channel `orm:"rel(fk)"`
}

func GenerateUserId(channelId, productId int, channelUserId string) (ret string) {
	return tool.Md5([]byte(fmt.Sprintf("%d%d%s", channelId, productId, channelUserId)))
}

func (this *LoginRequest) Init() {
	this.Id = -1
	this.ChannelId = -1
	this.ChannelUsername = ""
	this.ChannelUserid = ""
	this.Token = ""
	this.ProductId = -1
	this.IsDebug = false
	this.Userid = ""
}

func (m *LoginRequest) TableName() string {
	return "loginRequest"
}

func (m *LoginRequest) Insert() error {
	if _, err := orm.NewOrm().Insert(m); err != nil {
		return err
	}
	return nil
}

func (m *LoginRequest) Read(fields ...string) error {
	if err := orm.NewOrm().Read(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *LoginRequest) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *LoginRequest) Delete() error {
	if _, err := orm.NewOrm().Delete(m); err != nil {
		return err
	}
	return nil
}

func (m *LoginRequest) Query() orm.QuerySeter {
	return orm.NewOrm().QueryTable(m)
}
