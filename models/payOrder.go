package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type PayOrder struct {
	Id             int
	OrderId        string //订单号
	ChannelOrderId string
	PayAmount      int
	PayTime        time.Time `orm:"auto_now;type(datatime)"`
}

func (this *PayOrder) Init() {
	this.Id = -1
	this.OrderId = ""
	this.ChannelOrderId = ""
	this.PayAmount = -1
}

func (m *PayOrder) TableName() string {
	return "payOrder"
}

func (m *PayOrder) Insert() error {
	if _, err := orm.NewOrm().Insert(m); err != nil {
		return err
	}
	return nil
}

func (m *PayOrder) Read(fields ...string) error {
	if err := orm.NewOrm().Read(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *PayOrder) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *PayOrder) Delete() error {
	if _, err := orm.NewOrm().Delete(m); err != nil {
		return err
	}
	return nil
}

func (m *PayOrder) Query() orm.QuerySeter {
	return orm.NewOrm().QueryTable(m)
}
