package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type OrderRequest struct {
	Id             int
	OrderId        string    //订单id
	ChannelOrderId string    //渠道订单id
	UserId         string    //用户id
	ProductId      int       //产品id
	ChannelId      int       //渠道id
	ProductOrderId string    //产品订单号
	ReqAmount      int       //请求订单金额
	ReqTime        time.Time `orm:"auto_now_add;type(datatime)"` //请求时间
	State          int       //0初始 1完成渠道API通知回调 2完成产品通知API调用
	ProductCode    int       //产品通知API调用状态码  -1初始0成功1失败 ...
	ProductMessage string    //产品通知API调用消息
	ChannelLog     string    //渠道API通知回调日志
	CallbackUrl    string    //游戏服务器支付回调地址
}

func (this *OrderRequest) Init() {
	this.Id = -1
	this.OrderId = ""
	this.ChannelOrderId = ""
	this.UserId = ""
	this.ProductId = -1
	this.ChannelId = -1
	this.ProductOrderId = ""
	this.ReqAmount = -1
	this.State = 0
	this.ProductCode = -1
	this.ProductMessage = ""
	this.ChannelLog = ""
	this.CallbackUrl = ""
}

func (m *OrderRequest) TableName() string {
	return "orderRequest"
}

func (m *OrderRequest) Insert() error {
	if _, err := orm.NewOrm().Insert(m); err != nil {
		return err
	}
	return nil
}

func (m *OrderRequest) Read(fields ...string) error {
	if err := orm.NewOrm().Read(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *OrderRequest) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *OrderRequest) Delete() error {
	if _, err := orm.NewOrm().Delete(m); err != nil {
		return err
	}
	return nil
}

func (m *OrderRequest) Query() orm.QuerySeter {
	return orm.NewOrm().QueryTable(m)
}
