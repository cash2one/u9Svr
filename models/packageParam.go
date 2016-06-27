package models

import (
	"encoding/json"
	"github.com/astaxie/beego/orm"
	"time"
)

type PackageParam struct {
	Id          int
	ChannelId   int
	ProductId   int
	ProductName string
	PackageName string
	SignKeyFile string
	IconType    int8
	PackageIcon string
	JsonParam   string
	XmlParam    string
	ExtParam    string
	UpdateTime  time.Time `orm:"auto_now_add;type(datetime)"`
}

func (m *PackageParam) TableName() string {
	return "packageParam"
}

func (m *PackageParam) Insert() error {
	if _, err := orm.NewOrm().Insert(m); err != nil {
		return err
	}
	return nil
}

func (m *PackageParam) Read() error {
	if err := orm.NewOrm().Read(m); err != nil {
		return err
	}
	return nil
}

func (m *PackageParam) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *PackageParam) Delete() error {
	if _, err := orm.NewOrm().Delete(m); err != nil {
		return err
	}
	return nil
}

func (m *PackageParam) Query() orm.QuerySeter {
	return orm.NewOrm().QueryTable(m)
}

func GetChannelXmlParam(channelId, productId int) (ret map[string]string, err error) {
	m := new(PackageParam)

	err = m.Query().Filter("channelId", channelId).Filter("productId", productId).One(m)

	if err != nil {
		return
	}

	if err = json.Unmarshal([]byte(m.XmlParam), &ret); err != nil {
		return
	}
	return
}
