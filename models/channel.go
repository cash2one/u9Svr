package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type Channel struct {
	Id              int
	Name            string
	Type            string
	IsCustomPackage bool
	IsCustomSign    bool
	SdkVersion      string
	IconLeftTop     string
	IconLeftBottom  string
	IconRightTop    string
	IconRightBottom string
	UpdateTime      time.Time `orm:"auto_now;type(datetime)"`
	Desc            string
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
