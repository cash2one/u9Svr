package manager

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
	"u9/models"
	"u9/www/common"
)

type statisticPayItem struct {
	ChannelId   int
	ChannelName string
	ProductId   int
	ProductName string
	ReqAmount   string
	PayAmount   string
	OrderNum    int
	UserNum     int
	PayTime     time.Time `orm:"auto_now;type(datetime)"`
}

type productItem struct {
	Id   int
	Name string
}

type TitleInfo struct {
	ProductId   int
	ProductName string
	ChannelId   int
	ChannelName string
	StartTime   string
	EndTime     string
}

type StatPayCtl struct {
	BaseController
}
type Count struct {
	Count int64
}

func (this *StatPayCtl) Stat() {
	var err error
	var count Count

	var list []*statisticPayItem
	var productList []*productItem
	var channelList []*models.Channel
	var channel models.Channel

	var num int64
	var pageSize int64 = 15
	pageIndex, _ := this.GetInt64("page", 1)
	offset := (pageIndex - 1) * pageSize

	//查询数据
	queryBuilder, _ := orm.NewQueryBuilder("mysql")
	queryBuilder.Select("statChannelPay.channel_id", "statChannelPay.channel_name", "statChannelPay.product_id",
		"statChannelPay.product_name", "statChannelPay.req_amount", "statChannelPay.pay_amount",
		"statChannelPay.order_num", "statChannelPay.user_num", "statChannelPay.pay_time").
		From("statChannelPay").
		Where("statChannelPay.channel_id != 100").
		And("statChannelPay.product_id != 1000")
	//查询总数量
	numqueryBuilder, _ := orm.NewQueryBuilder("mysql")
	numqueryBuilder.Select("count(*) as count").From("statChannelPay").Where("statChannelPay.channel_id != 100").And("statChannelPay.product_id != 1000")
	//查询产品ID、名称
	ptQueryBuilder, _ := orm.NewQueryBuilder("mysql")
	ptQueryBuilder.Select("product.id", "product.name").
		From("product").
		OrderBy("id").Asc()
	//查询渠道ID、名称
	clCount, _ := channel.Query().Count()
	if clCount > 0 {
		channel.Query().OrderBy("id").All(&channelList)
	}

	if _, err = orm.NewOrm().Raw(ptQueryBuilder.String()).QueryRows(&productList); err != nil {
		beego.Error(err)
		return
	}

	var productId, channelId string
	var startDate, endDate string
	conditions := make(map[string]string, 0)

	pvProductId := this.Ctx.Request.PostForm["Product"]
	pvChannelId := this.Ctx.Request.PostForm["Channel"]
	pvStartDate := this.Ctx.Request.PostForm["StartDate"]
	pvEndDate := this.Ctx.Request.PostForm["EndDate"]
	if len(pvProductId) > 0 {
		productId = pvProductId[0]
	} else {
		productId = this.GetString("Product")
	}

	if len(pvChannelId) > 0 {
		channelId = pvChannelId[0]
	} else {
		channelId = this.GetString("Channel")
	}

	if len(pvStartDate) > 0 {
		startDate = pvStartDate[0]
	} else {
		startDate = this.GetString("StartDate")
	}

	if len(pvEndDate) > 0 {
		endDate = pvEndDate[0]
	} else {
		endDate = this.GetString("EndDate")
	}

	urlParamStr := ""

	if productId != "" {
		conditions["product_id"] = productId
		urlParamStr = urlParamStr + "&Product=" + productId
	}
	if channelId != "" {
		conditions["channel_id"] = channelId
		urlParamStr = urlParamStr + "&Channel=" + channelId
	}
	if startDate != "" {
		conditions["pay_time_start"] = "\"" + startDate + "\""
		urlParamStr = urlParamStr + "&StartDate=" + startDate
	}
	if endDate != "" {
		conditions["pay_time_end"] = "\"" + endDate + "\""
		urlParamStr = urlParamStr + "&EndDate=" + endDate
	}

	beego.Trace("productId:", productId, "#channelId:", channelId, "#startDate:", startDate, "#endDate:", endDate, "#urlParamStr:", urlParamStr)

	for k, v := range conditions {
		compareToken := "="
		if k == "pay_time_start" {
			k = "pay_time"
			compareToken = ">="
		} else if k == "pay_time_end" {
			k = "pay_time"
			compareToken = "<="
		}
		queryBuilder.And("statChannelPay." + k + compareToken + v)
		numqueryBuilder.And("statChannelPay." + k + compareToken + v)
		beego.Trace("k:", k, ";v:", v)
	}

	//查询总条数
	if err = orm.NewOrm().Raw(numqueryBuilder.String()).QueryRow(&count); err != nil {
		beego.Error(err)
		return
	}
	num = count.Count
	//查询数据
	queryBuilder.OrderBy("product_id", "channel_id", "pay_time").Desc().Limit(int(pageSize)).Offset(int(offset))
	if _, err = orm.NewOrm().Raw(queryBuilder.String()).QueryRows(&list); err != nil {
		beego.Error(err)
		return
	}

	if this.Ctx.Request.Method == "POST" {
		pageBar := common.NewPager(pageIndex, num, pageSize, "/manager/statistic/payList?page=%d"+urlParamStr).ToString()
		beego.Trace("urlParamStr:", urlParamStr)
		this.Data["pagebar"] = pageBar
		this.Data["payList"] = list
	} else {
		pageBar := common.NewPager(pageIndex, num, pageSize, "/manager/statistic/payList?page=%d"+urlParamStr).ToString()
		beego.Trace("urlParamStr:", urlParamStr)
		this.Data["pagebar"] = pageBar
		this.Data["payList"] = list
	}
	this.Data["productList"] = productList
	this.Data["channelList"] = channelList

	this.updateData()
	this.display()
}
