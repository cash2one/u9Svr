package cp

import (
	"errors"
	"github.com/astaxie/beego"
	"os"
	"time"
	"u9/models"
	"u9/www/android"
	"u9/www/common"
)

type PackageController struct {
	BaseController
	products             []*models.Product
	productVersions      []*models.ProductVersion
	channelPackageParams []*models.ChannelPackageParam
}

func (this *PackageController) List() {
	var err error
	pageBar := ""
	var pageSize int64 = 15
	var recordCount int64 = 0
	pageIndex, _ := this.GetInt64("page", 1)
	offset := (pageIndex - 1) * pageSize

	var packageTaskLists []*models.PackageTaskList
	var packageTaskList models.PackageTaskList

	defer func() {
		this.Data["list"] = packageTaskLists
		this.Data["pagebar"] = pageBar
		this.updateData()
		this.display()
	}()

	cpId := this.getCp().Id
	if recordCount, err = packageTaskList.Query().Filter("CpId", cpId).Count(); err != nil {
		beego.Error(err)
		return
	}
	urlPathFormat := "/cp/package/list?page=%d"
	pageBar = common.NewPager(pageIndex, recordCount, pageSize, urlPathFormat).ToString()
	if _, err = packageTaskList.Query().Filter("CpId", this.getCp().Id).
		OrderBy("-id").Limit(pageSize, offset).All(&packageTaskLists); err != nil {
		beego.Error(err)
		return
	}
}

// 根据产品ID号得到版本ID和渠道ID及相关参数
func (this *PackageController) GetVerAndClByPid() {
	type verParam struct {
		VerId   int    `json:"versionId"`
		VerCode string `json:"versionCode"`
	}

	type channelParam struct {
		PackageParamId int    `json:"packageParamId"`
		ChannelName    string `json:"channelName"`
	}

	type vcRetParam struct {
		VerParams []verParam     `json:"version"`
		XmlParams []channelParam `json:"channel"`
	}
	var ret vcRetParam

	defer func() {
		this.Data["json"] = &ret
		this.ServeJSON(true)
	}()

	productId, _ := this.GetInt("ProductId", 0)

	var productVersion models.ProductVersion
	productVersion.Query().Filter("product_id", productId).All(&this.productVersions)
	ret.VerParams = make([]verParam, len(this.productVersions))
	for i, v := range this.productVersions {

		ret.VerParams[i].VerId = v.Id
		ret.VerParams[i].VerCode = v.VersionCode
	}

	var channelPackageParam models.ChannelPackageParam
	channelPackageParam.Query().Filter("product_id", productId).All(&this.channelPackageParams)
	ret.XmlParams = make([]channelParam, len(this.channelPackageParams))
	for i, v := range this.channelPackageParams {
		ret.XmlParams[i].PackageParamId = v.Id
		ret.XmlParams[i].ChannelName = v.ChannelName
	}
}

func (this *PackageController) Add() {
	cp := this.getCp()
	if this.Ctx.Request.Method == "POST" {
		productId, _ := this.GetInt("Product", -1)
		versionId, _ := this.GetInt("Version", -1)
		packageParamId, _ := this.GetInt("Channel", -1)

		if productId == -1 || versionId == -1 || packageParamId == -1 {
			err := errors.New("productId,versionId,packageParamId is invalid.")
			this.showMsg(err.Error())
		}

		product := new(models.Product)
		product.Query().Filter("id", productId).One(product)

		packageTask := new(models.PackageTask)
		packageTask.CpId = cp.Id
		packageTask.PackageParamId = packageParamId
		packageTask.ProductVersionId = versionId
		packageTask.VersionUpdateTime = time.Now()
		packageTask.ChannelUpdateTime = time.Now()
		packageTask.State = 0

		if err := packageTask.Insert(); err != nil {
			this.showMsg(err.Error())
		}
		this.Redirect("/cp/package/list", 302)

	}

	var product models.Product
	product.Query().Filter("cp_id", cp.Id).OrderBy("id").All(&this.products)
	this.Data["products"] = this.products
	this.updateData()
	this.display()
}

func (this *PackageController) Package() {
	ret := "fail"
	defer func() {
		this.Data["json"] = &ret
		this.ServeJSON(true)
	}()

	id, _ := this.GetInt("Id")
	packageTask := models.PackageTask{Id: id}
	if err := packageTask.Read(); err != nil {
		return
	}
	if packageTask.State == 1 {
		ret = "wait"
		return
	}

	packageTask.State = 1
	if err := packageTask.Update("state"); err != nil {
		return
	}

	if publishApk, err := androidPackage.HandleTask(id); err != nil {
		packageTask.State = 3
		packageTask.Update("state")
	} else {
		if packageTask.PublishApk != "" {
			if err := os.Remove(packageTask.PublishApk); err != nil {
				beego.Warn(err)
			}
		}

		packageTask.State = 2
		packageTask.PublishApk = publishApk
		packageTask.PackageTime = time.Now()
		packageTask.Update("state", "publishApk", "packageTime")
		ret = "success"
	}
}

func (this *PackageController) Download() {
	id, _ := this.GetInt("Id")
	packageTask := models.PackageTask{Id: id}
	if err := packageTask.Read(); err != nil {
		this.Abort("404")
	} else {
		this.Ctx.Output.Download(packageTask.PublishApk)
	}
}

func (this *PackageController) Delete() {
	ret := "fail"
	defer func() {
		this.Data["json"] = &ret
		this.ServeJSON(true)
	}()

	id, _ := this.GetInt("Id")
	packageTask := models.PackageTask{Id: id}
	if err := packageTask.Read(); err != nil {
		return
	}
	if err := packageTask.Delete(); err != nil {
		return
	}

	ret = "success"
}
