package cp

import (
	//"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"strconv"
	"u9/models"
	"u9/www/common"
)

type proudctItem struct {
	Id          int
	CpId        int
	Direction   int //0 横屏  1 竖屏
	Name        string
	Code        string
	AppKey      string
	CallbackUrl string
	CpName      string
}

type ProductController struct {
	BaseController
}

func (this *ProductController) Edit() {

}

func (this *ProductController) Delete() {

}

func (this *ProductController) List() {
	pageBar := ""
	var productList []*proudctItem

	defer func() {
		this.Data["list"] = productList
		this.Data["pagebar"] = pageBar
		this.updateData()
		this.display()
	}()

	var err error
	var recordCount int64 = 0
	var product models.Product
	cpId := this.getCp().Id
	if recordCount, err = product.Query().Filter("CpId", cpId).Count(); err != nil {
		beego.Error(err)
		return
	}

	var pageSize int64 = 15
	pageIndex, _ := this.GetInt64("page", 1)
	offset := (pageIndex - 1) * pageSize

	queryBuilder, _ := orm.NewQueryBuilder("mysql")
	queryBuilder.Select("product.id", "product.direction", "product.code",
		"product.name", "product.app_key", "product.callback_url", "cp.cp_name").
		From("product").
		InnerJoin("cp").On("product.cp_id = cp.id").
		Where("product.cp_id=" + strconv.Itoa(cpId)).
		OrderBy("id").Asc().
		Limit(int(pageSize)).Offset(int(offset))

	if _, err = orm.NewOrm().Raw(queryBuilder.String()).QueryRows(&productList); err != nil {
		beego.Error(err)
		return
	}

	urlPathFormat := "/cp/product/list?page=%d"
	pageBar = common.NewPager(pageIndex, recordCount, pageSize, urlPathFormat).ToString()
}
