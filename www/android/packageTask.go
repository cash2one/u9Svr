package androidPackage

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"os"
	"u9/models"
	"u9/tool"
)

type PackageTaskHandle struct {
	packageTask    models.PackageTask
	product        models.Product
	productVersion models.ProductVersion
	channel        models.Channel
	packageParam   models.PackageParam

	apkName string

	res      *Res
	asset    *Asset
	smali    *Smali
	manifest *Manifest
	publish  *Publish
}

func HandleTask(packageTaskId int) (publishApk string, err error) {
	packageTaskHandle := new(PackageTaskHandle)

	defer func() {
		if rErr := recover(); rErr != nil {
			beego.Trace(rErr)
			msg := fmt.Sprintf("%v", rErr)
			err = errors.New(msg)
		}
	}()

	packageTaskHandle.Init(packageTaskId)
	publishApk = packageTaskHandle.Handle()
	return
}

func (this *PackageTaskHandle) Init(packageTaskId int) {
	this.packageTask = models.PackageTask{Id: packageTaskId}
	if err := this.packageTask.Read(); err != nil {
		panic(err)
	}

	this.productVersion = models.ProductVersion{Id: this.packageTask.ProductVersionId}
	if err := this.productVersion.Read(); err != nil {
		panic(err)
	}

	this.packageParam = models.PackageParam{Id: this.packageTask.PackageParamId}
	if err := this.packageParam.Read(); err != nil {
		panic(err)
	}

	this.channel = models.Channel{Id: this.packageParam.ChannelId}
	if err := this.channel.Read(); err != nil {
		panic(err)
	}

	this.product = models.Product{Id: this.packageParam.ProductId}
	if err := this.product.Read(); err != nil {
		panic(err)
	}

	this.apkName = GetApkName(&this.product, &this.productVersion)

	this.res = NewRes(this.packageTask.Id, &this.product, &this.productVersion, &this.packageParam)
	this.asset = NewAsset(this.packageTask.Id, &this.product, &this.productVersion, &this.packageParam)
	this.smali = NewSmali(this.packageTask.Id, &this.product, &this.productVersion)
	this.manifest = NewManifest(this.packageTask.Id,
		&this.product, &this.productVersion, &this.packageParam, &this.channel)
	this.publish = NewPublish(this.packageTask.Id, &this.product, &this.productVersion,
		&this.channel, &this.packageParam)
}

func (this *PackageTaskHandle) Prepare() {
	productPath := GetProductPath(&this.product, this.apkName)
	packageRootPath := GetPackagePath(this.packageTask.Id, "")
	beego.Trace("removeAll:", packageRootPath)
	if err := os.RemoveAll(packageRootPath); err != nil {
		panic(err)
	}
	beego.Trace("mkdirAll:", packageRootPath)
	if err := os.MkdirAll(packageRootPath, 0777); err != nil {
		panic(err)
	}

	beego.Trace("chmod:", packageRootPath)
	if err := os.Chmod(packageRootPath, 0777); err != nil {
		panic(err)
	}

	beego.Trace("copyDir:", productPath, packageRootPath)
	if err := tool.CopyDir(productPath, packageRootPath); err != nil {
		panic(err)
	}
}

func (this *PackageTaskHandle) Handle() (publishApk string) {
	this.Prepare()
	this.smali.Prepare()

	channelPath := GetChannelPath(&this.channel)
	packagePath := GetPackagePath(this.packageTask.Id, this.apkName)
	channelTar := channelPath + "/" + this.channel.Type + ".tar"
	beego.Trace("unTar:", packagePath, channelTar)
	if err := tool.UnTar(channelTar, packagePath); err != nil {
		panic(err)
	}

	beego.Trace("res handle...")
	this.res.Handle()

	beego.Trace("asset handle...")
	this.asset.Handle()

	beego.Trace("smali handle...")
	this.smali.Handle()

	beego.Trace("manifest handle...")
	this.manifest.Handle()

	beego.Trace("publishApk handle...")
	publishApk = this.publish.Handle()
	return
}
