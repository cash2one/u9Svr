package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type VStatUser struct {
	Id                  int64     `orm:"pk" json:"id"`
	ProductId           int       `json:"productId"`
	ProductName         string    `json:"productName"`
	ChannelId           int       `json:"channelId"`
	ChannelName         string    `json:"channelName"`
	LoginRequestIdCount int       `json:"loginRequestIdCount"`
	NewUserCount        int       `json:"newUserCount"`
	NewMobileCount      int       `json:"newMobileCount"`
	Time                time.Time `orm:"type(datetime)" json:"time"`
}

func (m *VStatUser) TableName() string {
	return "vStatUser"
}

func (m *VStatUser) Read() error {
	if err := orm.NewOrm().Read(m); err != nil {
		return err
	}
	return nil
}

func (m *VStatUser) Query() orm.QuerySeter {
	return orm.NewOrm().QueryTable(m)
}

func (m *VStatUser) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *VStatUser) Delete() error {
	if _, err := orm.NewOrm().Delete(m); err != nil {
		return err
	}
	return nil
}

func (m *VStatUser) Insert() error {
	if _, err := orm.NewOrm().Insert(m); err != nil {
		return err
	}
	return nil
}
