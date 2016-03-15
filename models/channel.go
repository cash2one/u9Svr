package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type Channel struct {
	Id              int
	Name            string `orm:"size(32)"`
	Type            string `orm:"size(16)"`
	IsCustomPackage bool
	IsCustomSign    bool
	SdkVersion      string    `orm:"size(32)"`
	IconLeftTop     string    `orm:"size(255)"`
	IconLeftBottom  string    `orm:"size(255)"`
	IconRightTop    string    `orm:"size(255)"`
	IconRightBottom string    `orm:"size(255)"`
	UpdateTime      time.Time `orm:"auto_now;type(datatime)"`
	Desc            string    `orm:"size(255)"`
}

func (m *Channel) TableName() string {
	return "channel"
}

func (m *Channel) Insert() error {
	if _, err := orm.NewOrm().Insert(m); err != nil {
		return err
	}
	return nil
}

func (m *Channel) Read(fields ...string) error {
	if err := orm.NewOrm().Read(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *Channel) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *Channel) Delete() error {
	if _, err := orm.NewOrm().Delete(m); err != nil {
		return err
	}
	return nil
}

func (m *Channel) Query() orm.QuerySeter {
	return orm.NewOrm().QueryTable(m)
}
