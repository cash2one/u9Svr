package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type PackageTaskView1 struct {
	Id                int
	ChannelName       string `orm:"size(32)"`
	ProductName       string `orm:"size(255)"`
	ProductCode       string `orm:"size(64)"`
	VersionCode       string `orm:"size(32)"`
	VersionName       string `orm:"size(32)"`
	ChannelType       string `orm:"size(16)"`
	ProductId         int
	CpId              int
	PackageTime       time.Time `orm:"type(datatime)"`
	VersionUpdateTime time.Time `orm:"type(datatime)"`
	ChannelUpdateTime time.Time `orm:"type(datatime)"`
	State             int       //0:初始 1正在打包 2:打包成功 3:打包失败
	PackageParamId    int
	ProductVersionId  int
}

func (m *PackageTaskView1) TableName() string {
	return "packageTaskView1"
}

func (m *PackageTaskView1) Insert() error {
	if _, err := orm.NewOrm().Insert(m); err != nil {
		return err
	}
	return nil
}

func (m *PackageTaskView1) Read() error {
	if err := orm.NewOrm().Read(m); err != nil {
		return err
	}
	return nil
}

func (m *PackageTaskView1) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *PackageTaskView1) Delete() error {
	if _, err := orm.NewOrm().Delete(m); err != nil {
		return err
	}
	return nil
}

func (m *PackageTaskView1) Query() orm.QuerySeter {
	return orm.NewOrm().QueryTable(m)
}
