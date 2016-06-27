package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type StatChannelPay struct {
	Id          string `orm:"pk"`
	ChannelId   int
	ChannelName string
	ProductId   int
	ProductName string
	ReqAmount   string
	PayAmount   string
	OrderNum    int
	UserNum     int
	PayTime     time.Time `orm:"auto_now;type(datetime)"`
}

func (m *StatChannelPay) TableName() string {
	return "statChannelPay"
}

func (m *StatChannelPay) Read() error {
	if err := orm.NewOrm().Read(m); err != nil {
		return err
	}
	return nil
}

func (m *StatChannelPay) Query() orm.QuerySeter {
	return orm.NewOrm().QueryTable(m)
}

func (m *StatChannelPay) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *StatChannelPay) Delete() error {
	if _, err := orm.NewOrm().Delete(m); err != nil {
		return err
	}
	return nil
}

func (m *StatChannelPay) Insert() error {
	if _, err := orm.NewOrm().Insert(m); err != nil {
		return err
	}
	return nil
}
