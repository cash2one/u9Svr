package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type ProductVersion struct {
	Id             int
	ProductId      int
	AppName        string
	VersionCode    string
	VersionName    string
	ApkUpdateTime  time.Time `orm:"type(datetime)"`
	ApkUpdateState int
	IconUrl        string
	UpdateTime     time.Time `orm:"auto_now_add;type(datetime)"`
}

func (m *ProductVersion) TableName() string {
	return "productVersion"
}

func (m *ProductVersion) Insert() error {
	if _, err := orm.NewOrm().Insert(m); err != nil {
		return err
	}
	return nil
}

func (m *ProductVersion) Read() error {
	if err := orm.NewOrm().Read(m); err != nil {
		return err
	}
	return nil
}

func (m *ProductVersion) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *ProductVersion) Delete() error {
	if _, err := orm.NewOrm().Delete(m); err != nil {
		return err
	}
	return nil
}

func (m *ProductVersion) Query() orm.QuerySeter {
	return orm.NewOrm().QueryTable(m)
}
