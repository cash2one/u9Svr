package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type PackageTaskList struct {
	Id                int
	ChannelName       string
	ChannelType       string
	ProductName       string
	ProductCode       string
	VersionCode       string
	VersionName       string
	ProductId         int
	CpId              int
	PackageTime       time.Time `orm:"type(datetime)"`
	VersionUpdateTime time.Time `orm:"type(datetime)"`
	ChannelUpdateTime time.Time `orm:"type(datetime)"`
	State             int       //0:初始 1正在打包 2:打包成功 3:打包失败
	PackageParamId    int
	ProductVersionId  int
}

func (m *PackageTaskList) TableName() string {
	return "packageTaskList"
}

func (m *PackageTaskList) Read() error {
	if err := orm.NewOrm().Read(m); err != nil {
		return err
	}
	return nil
}

func (m *PackageTaskList) Query() orm.QuerySeter {
	return orm.NewOrm().QueryTable(m)
}

func (m *PackageTaskList) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *PackageTaskList) Delete() error {
	if _, err := orm.NewOrm().Delete(m); err != nil {
		return err
	}
	return nil
}

func (m *PackageTaskList) Insert() error {
	if _, err := orm.NewOrm().Insert(m); err != nil {
		return err
	}
	return nil
}
