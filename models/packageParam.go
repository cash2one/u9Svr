package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

type PackageParam struct {
	Id          int
	ChannelId   int
	ProductId   int
	ProductName string `orm:"size(255)"`
	PackageName string `orm:"size(255)"`
	SignKeyFile string `orm:"size(255)"`
	IconType    int8
	PackageIcon string    `orm:"size(512)"`
	JsonParam   string    `orm:"size(512)"`
	XmlParam    string    `orm:"size(2048)"`
	ExtParam    string    `orm:"size(512)"`
	UpdateTime  time.Time `orm:"auto_now_add;type(datatime)"`
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

func (m *PackageParam) GetXmlParam(channelId, productId int, key string) (ret string, err error) {
	err = m.Query().Filter("channelId", channelId).Filter("productId", productId).One(m)
	if err != nil {
		return
	}

	args := new(map[string]string)
	if err = json.Unmarshal([]byte(m.XmlParam), args); err != nil {
		beego.Error(m.XmlParam)
		return
	}

	ok := false
	if ret, ok = (*args)[key]; !ok {
		msg := fmt.Sprintf("PackageParam %s is empty.", key)
		err = errors.New(msg)
		return
	}
	return
}
