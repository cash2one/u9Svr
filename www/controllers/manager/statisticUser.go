package manager

import (
	"database/sql"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"math"
	"strconv"
	"time"
	"u9/models"
)

type StatUserCtl struct {
	BaseController
}

func (this *StatUserCtl) URLMapping() {
	this.Mapping("Query", this.Query)
	this.Mapping("Product", this.Product)
	this.Mapping("Channel", this.Channel)
}

func (this *StatUserCtl) Get() {
	// var err error

	// defer func() {
	// 	if err != nil {
	// 		beego.Error(err)
	// 	}
	// }()

	this.Layout = "manager/layout2.html"
	this.TplName = "manager/user_statistic.html"
	this.updateData()
	// var productList []*models.Product
	// if _, err = new(models.Product).Query().OrderBy("id").All(&productList); err != nil {
	// 	return
	// }
	// this.Data["productList"] = productList

	// var channelList []*models.Channel
	// if _, err = new(models.Channel).Query().OrderBy("id").All(&channelList); err != nil {
	// 	return
	// }
	// this.Data["channelList"] = channelList
}

func (this *StatUserCtl) Update(productId, channelId int, startTime, endTime string) (num int64, err error) {
	format := "CALL statUser(%d, %d,'%s','%s')"
	sqlText := fmt.Sprintf(format, productId, channelId, startTime, endTime)

	var rawPreparer orm.RawPreparer
	if rawPreparer, err = orm.NewOrm().Raw(sqlText).Prepare(); err != nil {
		return
	}
	defer func() {
		rawPreparer.Close()
	}()

	var result sql.Result
	if result, err = rawPreparer.Exec(); err != nil {
		return
	}

	if num, err = result.RowsAffected(); err != nil {
		return
	}

	return
}

// @router /manager/statistic/user/query/ [post]
func (this *StatUserCtl) Query() {
	const dateFormat = "2006-01-02"

	type userStaticData struct {
		TotalPages int64               `json:"totalPages"`
		Data       []*models.VStatUser `json:"data"`
	}
	var ret userStaticData

	defer func() {
		this.Data["json"] = &ret
		this.ServeJSON(true)
	}()

	productId, _ := this.GetInt("product", 0)
	channelId, _ := this.GetInt("channel", 0)

	var err error
	strStartDate := this.GetString("startDate")
	var startDate time.Time
	if startDate, err = time.Parse(dateFormat, strStartDate); err == nil {
		strStartDate = startDate.Format(dateFormat)
	} else {
		strStartDate = ""
		//beego.Warn(err)
	}

	strEndDate := this.GetString("endDate")
	var endDate time.Time
	if endDate, err = time.Parse(dateFormat, strEndDate); err == nil {
		strEndDate = endDate.AddDate(0, 0, 1).Format(dateFormat)
	} else {
		strEndDate = ""
		//beego.Warn(err)
	}

	if this.GetString("updateData") == "true" {
		if _, err = this.Update(productId, channelId, strStartDate, strEndDate); err != nil {
			beego.Error(err)
			return
		}
	}

	var statUserList models.VStatUser
	qs := statUserList.Query()
	if productId > 0 {
		qs = qs.Filter("productId", productId)
	} else {
		qs = qs.Exclude("productId", 1000)
	}
	if channelId > 0 {
		qs = qs.Filter("channelId", channelId)
	} else {
		qs = qs.Exclude("channelId", 100)
	}
	if strStartDate != "" {
		qs = qs.Filter("time__gte", strStartDate)
	}
	if strEndDate != "" {
		qs = qs.Filter("time__lt", strEndDate)
	}

	page, _ := this.GetInt64("page", 1)
	pageSize, _ := this.GetInt64("pageSize", 15)
	offset := (page - 1) * pageSize

	var recordCount int64
	if recordCount, err = qs.Count(); err != nil {
		beego.Error(err)
		return
	}

	ret.TotalPages = int64(math.Ceil(float64(recordCount) / float64(pageSize)))
	if _, err = qs.OrderBy("product_id", "channel_id", "time").Limit(pageSize).Offset(offset).All(&ret.Data); err != nil {
		beego.Error(err)
		return
	}
}

// @router /manager/statistic/user/product/ [post]
func (this *StatUserCtl) Product() {
	type productRet struct {
		Data []struct {
			ProductId   int    `json:"id"`
			ProductName string `json:"name"`
		} `json:"data"`
	}
	var ret productRet

	defer func() {
		this.Data["json"] = &ret
		this.ServeJSON(true)
	}()

	sql := "SELECT DISTINCT product_id, product_name FROM vProductChannel"
	sql = sql + " WHERE channel_id != 100 AND product_id != 1000"
	sql = sql + " ORDER BY product_id"
	if _, err := orm.NewOrm().Raw(sql).QueryRows(&ret.Data); err != nil {
		beego.Error(err)
	}
}

// @router /manager/statistic/user/channel/ [post]
func (this *StatUserCtl) Channel() {
	type channelRet struct {
		Data []struct {
			ChannelId   int    `json:"id"`
			ChannelName string `json:"name"`
		} `json:"data"`
	}
	var ret channelRet

	defer func() {
		this.Data["json"] = &ret
		this.ServeJSON(true)
	}()

	sql := "SELECT DISTINCT channel_id, channel_name FROM vProductChannel"
	sql = sql + " WHERE channel_id != 100 AND product_id != 1000"
	productId, _ := this.GetInt("product", 0)
	if productId > 0 {
		sql = sql + " AND product_id = " + strconv.Itoa(productId)
	}
	sql = sql + " ORDER BY channel_id"
	if _, err := orm.NewOrm().Raw(sql).QueryRows(&ret.Data); err != nil {
		beego.Error(err)
	}
}
