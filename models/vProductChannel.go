package models

import (
	"github.com/astaxie/beego/orm"
)

type VProductChannel struct {
	Id          int64  `orm:"pk" json:"id"`
	ChannelId   int    `json:"channelId"`
	ChannelName string `json:"channelName"`
	ProductId   int    `json:"productId"`
	ProductName string `json:"productName"`
}

func (m *VProductChannel) TableName() string {
	return "vProductChannel"
}

func (m *VProductChannel) Read() error {
	if err := orm.NewOrm().Read(m); err != nil {
		return err
	}
	return nil
}

func (m *VProductChannel) Query() orm.QuerySeter {
	return orm.NewOrm().QueryTable(m)
}

func (m *VProductChannel) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *VProductChannel) Delete() error {
	if _, err := orm.NewOrm().Delete(m); err != nil {
		return err
	}
	return nil
}

func (m *VProductChannel) Insert() error {
	if _, err := orm.NewOrm().Insert(m); err != nil {
		return err
	}
	return nil
}
