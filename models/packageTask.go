package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type PackageTask struct {
	Id                int
	CpId              int
	PackageParamId    int
	ProductVersionId  int
	PackageTime       time.Time `orm:"type(datatime)"`
	VersionUpdateTime time.Time `orm:"type(datatime)"`
	ChannelUpdateTime time.Time `orm:"type(datatime)"`
	State             int       //0:初始 1正在打包 2:打包成功 3:打包失败'
	PublishApk        string    `orm:"size(255)"`
	Log               string    `orm:"size(512)"`
}

func (m *PackageTask) TableName() string {
	return "packageTask"
}

func (m *PackageTask) Insert() error {
	if _, err := orm.NewOrm().Insert(m); err != nil {
		return err
	}
	return nil
}

func (m *PackageTask) Read() error {
	if err := orm.NewOrm().Read(m); err != nil {
		return err
	}
	return nil
}

func (m *PackageTask) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *PackageTask) Delete() error {
	if _, err := orm.NewOrm().Delete(m); err != nil {
		return err
	}
	return nil
}

func (m *PackageTask) Query() orm.QuerySeter {
	return orm.NewOrm().QueryTable(m)
}
