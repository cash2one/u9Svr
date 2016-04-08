package models

import (
	// "fmt"
	"github.com/astaxie/beego/orm"
	"time"
)

type LoginRequestLog struct {
	Id             int
	LoginRequestId int
	LoginTime      time.Time `orm:"auto_now_add;type(datatime)"`
}

func (m *LoginRequestLog) TableName() string {
	return "loginRequestLog"
}

func (m *LoginRequestLog) Insert() error {
	if _, err := orm.NewOrm().Insert(m); err != nil {
		return err
	}
	return nil
}

func (m *LoginRequestLog) Read(fields ...string) error {
	if err := orm.NewOrm().Read(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *LoginRequestLog) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *LoginRequestLog) Delete() error {
	if _, err := orm.NewOrm().Delete(m); err != nil {
		return err
	}
	return nil
}

func (m *LoginRequestLog) Query() orm.QuerySeter {
	return orm.NewOrm().QueryTable(m)
}
