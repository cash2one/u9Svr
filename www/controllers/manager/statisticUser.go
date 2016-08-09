package manager

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/tealeg/xlsx"
	"math"
	"net/http"
	"strconv"
	"time"
	"u9/models"
)

type StatUserCtl struct {
	BaseController
}

type condition struct {
	ProductId  int
	ChannelId  int
	StartDate  string
	EndDate    string
	UpdateData bool
	Page       int
	PageSize   int
}

type userStaticData struct {
	Count               int64               `json:"count"`
	TotalPages          int                 `json:"totalPages"`
	NewUserCount        int                 `json:"newUserCount"`
	NewMobileCount      int                 `json:"newMobileCount"`
	LoginRequestIdCount int                 `json:"loginRequestIdCount"`
	Data                []*models.VStatUser `json:"data"`
}

func (this *StatUserCtl) URLMapping() {
	this.Mapping("CommonQuery", this.CommonQuery)
	this.Mapping("GroupQuery", this.GroupQuery)
	this.Mapping("Product", this.Product)
	this.Mapping("Channel", this.Channel)
	this.Mapping("Export", this.Export)
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

func (this *StatUserCtl) update(cond *condition) (num int64, err error) {
	if !cond.UpdateData {
		return
	}
	const format = "CALL statUser(%d, %d,'%s','%s')"
	sqlText := fmt.Sprintf(format,
		cond.ProductId, cond.ChannelId,
		cond.StartDate, cond.EndDate)

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

func (this *StatUserCtl) getCond() (ret *condition) {
	const dateFormat = "2006-01-02"

	ret = new(condition)
	ret.ProductId, _ = this.GetInt("product", 0)
	ret.ChannelId, _ = this.GetInt("channel", 0)

	var err error
	ret.StartDate = this.GetString("startDate")
	var startDate time.Time
	if startDate, err = time.Parse(dateFormat, ret.StartDate); err == nil {
		ret.StartDate = startDate.Format(dateFormat)
	} else {
		ret.StartDate = ""
	}

	ret.EndDate = this.GetString("endDate")
	var endDate time.Time
	if endDate, err = time.Parse(dateFormat, ret.EndDate); err == nil {
		ret.EndDate = endDate.AddDate(0, 0, 1).Format(dateFormat)
	} else {
		ret.EndDate = ""
	}

	ret.UpdateData = this.GetString("updateData") == "true"

	ret.Page, _ = this.GetInt("page", 0)
	ret.PageSize, _ = this.GetInt("pageSize", 15)
	return
}

func (this *StatUserCtl) commonQuery(cond *condition) (ret *userStaticData, err error) {
	const dateFormat = "2006-01-02"

	ret = &userStaticData{
		TotalPages:          0,
		NewUserCount:        -1,
		NewMobileCount:      -1,
		LoginRequestIdCount: -1}

	var statUserList models.VStatUser
	qs := statUserList.Query()
	if cond.ProductId > 0 {
		qs = qs.Filter("productId", cond.ProductId)
	} else {
		qs = qs.Exclude("productId", 1000)
	}
	if cond.ChannelId > 0 {
		qs = qs.Filter("channelId", cond.ChannelId)
	} else {
		qs = qs.Exclude("channelId", 100)
	}
	if cond.StartDate != "" {
		qs = qs.Filter("time__gte", cond.StartDate)
	}
	if cond.EndDate != "" {
		qs = qs.Filter("time__lt", cond.EndDate)
	}

	qs = qs.OrderBy("product_id", "channel_id", "time")

	switch cond.Page {
	case 0:
		ret.TotalPages = 1
	case 1:
		sql := "SELECT Sum(vStatUser.new_user_count) AS new_user_count,"
		sql = sql + " Sum(vStatUser.new_mobile_count) AS new_mobile_count,"
		sql = sql + " Sum(vStatUser.login_request_id_count) AS login_request_id_count"
		sql = sql + " FROM vStatUser WHERE product_id != 1000 AND channel_id != 100"
		if cond.ProductId > 0 {
			sql = sql + " AND product_id =" + strconv.Itoa(cond.ProductId)
		}
		if cond.ChannelId > 0 {
			sql = sql + " AND channel_id =" + strconv.Itoa(cond.ChannelId)
		}
		if cond.StartDate != "" {
			sql = sql + " AND time >=" + `"` + cond.StartDate + `"`
		}
		if cond.EndDate != "" {
			sql = sql + " AND time <" + `"` + cond.EndDate + `"`
		}
		//beego.Trace(sql)
		if err = orm.NewOrm().Raw(sql).QueryRow(&ret); err != nil {
			beego.Error(err)
			beego.Error(sql)
			return
		}
		//beego.Trace(fmt.Sprintf("%+v", ret))
		fallthrough
	default:
		offset := int64((cond.Page - 1) * cond.PageSize)

		if ret.Count, err = qs.Count(); err != nil {
			beego.Error(err)
			return
		}
		ret.TotalPages = int(math.Ceil(float64(ret.Count) / float64(cond.PageSize)))
		qs = qs.Limit(cond.PageSize).Offset(offset)
	}

	if _, err = qs.All(&ret.Data); err != nil {
		beego.Error(err)
		return
	}

	return
}

// @router /manager/statistic/user/commonQuery/ [post]
func (this *StatUserCtl) CommonQuery() {
	var err error
	ret := &userStaticData{}

	cond := this.getCond()

	defer func() {
		this.Data["json"] = ret
		//beego.Trace(fmt.Sprintf("%+v", cond))
		//beego.Trace(fmt.Sprintf("%+v", ret))
		this.ServeJSON(true)
	}()

	if _, err = this.update(cond); err != nil {
		beego.Error(err)
		return
	}

	if ret, err = this.commonQuery(cond); err != nil {
		beego.Error(err)
		return
	}
}

func (this *StatUserCtl) groupQuery(cond *condition) (ret *userStaticData, err error) {
	const dateFormat = "2006-01-02"

	ret = &userStaticData{
		TotalPages:          0,
		NewUserCount:        -1,
		NewMobileCount:      -1,
		LoginRequestIdCount: -1}

	sql := "SELECT product_id, product_name, channel_id, channel_name,"
	sql = sql + "Sum(login_request_id_count) AS login_request_id_count,"
	sql = sql + "Sum(new_user_count) AS new_user_count,"
	sql = sql + "Sum(new_mobile_count) AS new_mobile_count"
	sql = sql + " FROM vStatUser WHERE product_id != 1000 AND channel_id != 100"
	if cond.ProductId > 0 {
		sql = sql + " AND product_id =" + strconv.Itoa(cond.ProductId)
	}
	if cond.ChannelId > 0 {
		sql = sql + " AND channel_id =" + strconv.Itoa(cond.ChannelId)
	}
	if cond.StartDate != "" {
		sql = sql + " AND time >=" + `"` + cond.StartDate + `"`
	}
	if cond.EndDate != "" {
		sql = sql + " AND time <" + `"` + cond.EndDate + `"`
	}
	sql = sql + " GROUP BY product_id,channel_id"
	sql = sql + " ORDER BY product_id,channel_id"
	//beego.Trace(sql)
	switch cond.Page {
	case 0:
		ret.TotalPages = 1
	case 1:
		groupSql := "SELECT SUM(login_request_id_count) AS login_request_id_count,"
		groupSql = groupSql + " SUM(new_user_count) AS new_user_count,"
		groupSql = groupSql + " SUM(new_mobile_count) AS new_mobile_count"
		groupSql = groupSql + " FROM (" + sql + ") AS vStatUser"
		if err = orm.NewOrm().Raw(groupSql).QueryRow(&ret); err != nil {
			beego.Error(err)
			beego.Error(groupSql)
			return
		}
		//beego.Trace(groupSql)
		fallthrough
	default:
		countSql := "SELECT COUNT(*) AS count FROM (" + sql + ") AS vStatUser"
		if err = orm.NewOrm().Raw(countSql).QueryRow(&ret); err != nil {
			beego.Error(err)
			beego.Error(countSql)
			return
		}
		//beego.Trace(countSql)
		ret.TotalPages = int(math.Ceil(float64(ret.Count) / float64(cond.PageSize)))
		sql = sql + " LIMIT " + strconv.Itoa(cond.PageSize)
		offset := (cond.Page - 1) * cond.PageSize
		sql = sql + " OFFSET " + strconv.Itoa(offset)
	}

	if _, err = orm.NewOrm().Raw(sql).QueryRows(&ret.Data); err != nil {
		beego.Error(err)
		beego.Error(sql)
		return
	}
	return
}

// @router /manager/statistic/user/groupQuery/ [post]
func (this *StatUserCtl) GroupQuery() {
	var err error
	ret := &userStaticData{}

	cond := this.getCond()

	defer func() {
		this.Data["json"] = ret
		//beego.Trace(fmt.Sprintf("%+v", cond))
		//beego.Trace(fmt.Sprintf("%+v", ret))
		this.ServeJSON(true)
	}()

	if _, err = this.update(cond); err != nil {
		beego.Error(err)
		return
	}

	if ret, err = this.groupQuery(cond); err != nil {
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

// @router /manager/statistic/user/export/ [*]
func (this *StatUserCtl) Export() {
	var err error

	defer func() {
		if err != nil {
			this.Abort("404")
		}
	}()

	cond := this.getCond()
	file := xlsx.NewFile()
	if _, err = this.commonExport(cond, file); err != nil {
		return
	}

	if _, err = this.groupExport(cond, file); err != nil {
		return
	}

	if err = this.download(file, "用户统计.xlsx"); err != nil {
		beego.Error(err)
		return
	}
}

func (this *StatUserCtl) commonExport(cond *condition, file *xlsx.File) (sheet *xlsx.Sheet, err error) {
	titles := [...]string{"游戏", "渠道", "登录用户", "新增用户", "激活设备", "日期"}

	ret := &userStaticData{}
	if ret, err = this.commonQuery(cond); err != nil {
		beego.Error(err)
		return
	}

	style := xlsx.NewStyle()

	alignment := *xlsx.DefaultAlignment()
	alignment.Horizontal = "center"
	style.Alignment = alignment
	style.ApplyAlignment = true

	font := *xlsx.DefaultFont()
	font.Bold = true
	font.Color = "00FFFFFF"
	style.Font = font
	style.ApplyFont = true

	fill := *xlsx.NewFill("solid", "00082E54", "00000000")
	style.Fill = fill
	style.ApplyFill = true

	if sheet, err = file.AddSheet("普通统计"); err != nil {
		beego.Error(err)
		return
	}

	var cell *xlsx.Cell
	row := sheet.AddRow()
	for _, item := range titles {
		cell = row.AddCell()
		cell.Value = item
		cell.SetStyle(style)
	}

	for _, item := range ret.Data {
		row = sheet.AddRow()
		cell = row.AddCell()
		cell.SetString(item.ProductName)
		cell = row.AddCell()
		cell.SetString(item.ChannelName)
		cell = row.AddCell()
		cell.SetInt(item.LoginRequestIdCount)
		cell = row.AddCell()
		cell.SetInt(item.NewUserCount)
		cell = row.AddCell()
		cell.SetInt(item.NewMobileCount)
		cell = row.AddCell()
		cell.SetDate(item.Time)
	}

	row = sheet.AddRow()
	cell = row.AddCell()
	cell.Value = "汇总"
	cell = row.AddCell()

	cell = row.AddCell()
	length := len(ret.Data) + 1
	format := "SUM(%s%d:%s%d)"
	cell.SetFormula(fmt.Sprintf(format, "C", 2, "C", length))
	cell = row.AddCell()
	cell.SetFormula(fmt.Sprintf(format, "D", 2, "D", length))
	cell = row.AddCell()
	cell.SetFormula(fmt.Sprintf(format, "E", 2, "E", length))

	return
}

func (this *StatUserCtl) groupExport(cond *condition, file *xlsx.File) (sheet *xlsx.Sheet, err error) {
	titles := [...]string{"游戏", "渠道", "登录用户", "新增用户", "激活设备"}

	ret := &userStaticData{}
	if ret, err = this.groupQuery(cond); err != nil {
		beego.Error(err)
		return
	}

	style := xlsx.NewStyle()

	alignment := *xlsx.DefaultAlignment()
	alignment.Horizontal = "center"
	style.Alignment = alignment
	style.ApplyAlignment = true

	font := *xlsx.DefaultFont()
	font.Bold = true
	font.Color = "00FFFFFF"
	style.Font = font
	style.ApplyFont = true

	fill := *xlsx.NewFill("solid", "00082E54", "00000000")
	style.Fill = fill
	style.ApplyFill = true

	if sheet, err = file.AddSheet("分类统计"); err != nil {
		beego.Error(err)
		return
	}

	var cell *xlsx.Cell
	row := sheet.AddRow()
	for _, item := range titles {
		cell := row.AddCell()
		cell.Value = item
		cell.SetStyle(style)
	}

	for _, item := range ret.Data {
		row = sheet.AddRow()
		cell = row.AddCell()
		cell.SetString(item.ProductName)
		cell = row.AddCell()
		cell.SetString(item.ChannelName)
		cell = row.AddCell()
		cell.SetInt(item.LoginRequestIdCount)
		cell = row.AddCell()
		cell.SetInt(item.NewUserCount)
		cell = row.AddCell()
		cell.SetInt(item.NewMobileCount)
	}

	row = sheet.AddRow()
	cell = row.AddCell()
	cell.Value = "汇总"
	cell = row.AddCell()

	cell = row.AddCell()
	length := len(ret.Data) + 1
	format := "SUM(%s%d:%s%d)"
	cell.SetFormula(fmt.Sprintf(format, "C", 2, "C", length))
	cell = row.AddCell()
	cell.SetFormula(fmt.Sprintf(format, "D", 2, "D", length))
	cell = row.AddCell()
	cell.SetFormula(fmt.Sprintf(format, "E", 2, "E", length))

	return
}

func (this *StatUserCtl) download(file *xlsx.File, filename string) (err error) {
	output := this.Ctx.Output
	output.Header("Content-Description", "File Transfer")
	output.Header("Content-Type", "application/octet-stream")
	output.Header("Content-Disposition", "attachment; filename="+filename)
	output.Header("Content-Transfer-Encoding", "binary")
	output.Header("Expires", "0")
	output.Header("Cache-Control", "must-revalidate")
	output.Header("Pragma", "public")

	buf := bytes.NewBuffer([]byte(""))
	if err = file.Write(buf); err != nil {
		beego.Error(err)
		return
	}
	content := bytes.NewReader(buf.Bytes())
	http.ServeContent(output.Context.ResponseWriter, output.Context.Request, filename, time.Now(), content)
	return
}
