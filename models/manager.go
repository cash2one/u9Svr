package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type Manager struct {
	Id            int
	Name          string    `orm:"unique;size(15)"`
	Password      string    `orm:"size(32)"`
	Email         string    `orm:"size(50)"`
	LastLoginTime time.Time `orm:"auto_now_add;type(datetime)"`
	LastLoginIp   string    `orm:"size(32)"`
	LoginCount    int
	State         int8
}

func (m *Manager) TableName() string {
	return "manager"
}

func (m *Manager) Insert() error {
	if _, err := orm.NewOrm().Insert(m); err != nil {
		return err
	}
	return nil
}

func (m *Manager) Read(fields ...string) error {
	if err := orm.NewOrm().Read(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *Manager) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *Manager) Delete() error {
	if _, err := orm.NewOrm().Delete(m); err != nil {
		return err
	}
	return nil
}

func (m *Manager) Query() orm.QuerySeter {
	return orm.NewOrm().QueryTable(m)
}
