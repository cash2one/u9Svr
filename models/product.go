package models

import (
	"github.com/astaxie/beego/orm"
)

type Product struct {
	Id          int
	CpId        int
	Direction   int //0 横屏  1 竖屏
	Name        string
	Code        string
	AppKey      string
	CallbackUrl string
}

func (this *Product) Init() {
	this.Id = -1
	this.CpId = -1
	this.Direction = -1
	this.Name = ""
	this.Code = ""
	this.AppKey = ""
	this.CallbackUrl = ""
}

func (m *Product) Insert() error {
	if _, err := orm.NewOrm().Insert(m); err != nil {
		return err
	}
	return nil
}

func (m *Product) Read() error {
	if err := orm.NewOrm().Read(m); err != nil {
		return err
	}
	return nil
}

func (m *Product) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *Product) Delete() error {
	if _, err := orm.NewOrm().Delete(m); err != nil {
		return err
	}
	return nil
}

func (m *Product) Query() orm.QuerySeter {
	return orm.NewOrm().QueryTable(m)
}
