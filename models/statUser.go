package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type StatUser struct {
	Id                  int64     `orm:"pk" json:"id"`
	ChannelId           int       `json:"channelId"`
	ProductId           int       `json:"productId"`
	LoginRequestIdCount int       `json:"loginRequestIdCount"`
	NewUserCount        int       `json:"newUserCount"`
	NewMobileCount      int       `json:"newMobileCount"`
	Time                time.Time `orm:"type(datetime)" json:"time"`
}

func (m *StatUser) TableName() string {
	return "statUser"
}

func (m *StatUser) Read() error {
	if err := orm.NewOrm().Read(m); err != nil {
		return err
	}
	return nil
}

func (m *StatUser) Query() orm.QuerySeter {
	return orm.NewOrm().QueryTable(m)
}

func (m *StatUser) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *StatUser) Delete() error {
	if _, err := orm.NewOrm().Delete(m); err != nil {
		return err
	}
	return nil
}

func (m *StatUser) Insert() error {
	if _, err := orm.NewOrm().Insert(m); err != nil {
		return err
	}
	return nil
}
