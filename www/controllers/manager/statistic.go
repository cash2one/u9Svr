package manager

import (
	"github.com/astaxie/beego"
	// "github.com/astaxie/beego/validation"
	// "strings"
	"github.com/astaxie/beego/orm"
	"u9/models"
	"u9/www/common"
)

type StatisticController struct {
	BaseController
}

type payStatisticItem struct {
	ProductId int
}

//支付统计列表
func (this *StatisticController) PayList() {
	var page int64
	var pagesize int64 = 10
	var list []*models.StatChannelPay
	var statChannelPay models.StatChannelPay
	var payStatisticList []*payStatisticItem

	if page, _ = this.GetInt64("page"); page < 1 {
		page = 1
	}
	offset := (page - 1) * pagesize

	queryBuilder, _ := orm.NewQueryBuilder("mysql")
	queryBuilder.Select("product_id").
		From("statChannelPay").
		Limit(int(pagesize)).Offset(int(offset))

	beego.Trace(queryBuilder.String())

	if _, err := orm.NewOrm().Raw(queryBuilder.String()).QueryRows(&payStatisticList); err != nil {
		beego.Error(err)
		return
	}
	for _, v := range payStatisticList {
		beego.Trace(v)
	}

	count, err := statChannelPay.Query().Count()
	beego.Trace(count)
	beego.Trace(err)
	if count > 0 {
		count, err = statChannelPay.Query().OrderBy("-channelId").Limit(pagesize, offset).All(&list)
	}
	beego.Trace(count)
	beego.Trace(err)

	for _, v := range list {
		beego.Trace(v)
	}
	this.Data["payList"] = list
	pageBar := common.NewPager(page, count, pagesize, "/manager/statistic/payList?page=%d").ToString()
	this.Data["pagebar"] = pageBar

	this.updateData()
	this.display()
}
