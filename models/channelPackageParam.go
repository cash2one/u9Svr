package models

import (
	"github.com/astaxie/beego/orm"
)

type ChannelPackageParam struct {
	Id          int
	ProductId   int
	ChannelId   string
	ChannelName string
}

func (m *ChannelPackageParam) TableName() string {
	return "channelPackageParam"
}

func (m *ChannelPackageParam) Read(fields ...string) error {
	if err := orm.NewOrm().Read(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *ChannelPackageParam) Query() orm.QuerySeter {
	return orm.NewOrm().QueryTable(m)
}
