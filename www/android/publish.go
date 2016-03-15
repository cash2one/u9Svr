package androidPackage

import (
	"github.com/astaxie/beego"
	"u9/models"
	"u9/tool"
	"u9/tool/android"
)

const (
	defaultSignKeyFile = "package/key/hygame.keystore"
	signPassWord       = "hygame147"
	signAliasName      = "hygame"
)

type Publish struct {
	packageTaskId  int
	product        *models.Product
	productVersion *models.ProductVersion
	channel        *models.Channel
	packageParam   *models.PackageParam
}

func NewPublish(packageTaskId int, product *models.Product, productVersion *models.ProductVersion,
	channel *models.Channel, packageParam *models.PackageParam) *Publish {
	ret := new(Publish)
	ret.packageTaskId, ret.product, ret.productVersion, ret.channel, ret.packageParam =
		packageTaskId, product, productVersion, channel, packageParam
	return ret
}

func (this *Publish) Handle() (publishApk string) {
	//回编译
	beego.Trace("CompileApk")
	apkName := GetApkName(this.product, this.productVersion)
	packageRootPath := GetPackagePath(this.packageTaskId, "")
	compilePath := packageRootPath + "/" + apkName
	compileApk := compilePath + ".apk"
	if err := android.CompileApk(compilePath, compileApk); err != nil {
		panic(err)
	}

	//创建正式包文件夹
	beego.Trace("CreateDir")
	publishPath := GetPublishPath(this.product, this.productVersion)
	if err := tool.CreateDir(publishPath); err != nil {
		panic(err)
	}

	//签名
	beego.Trace("ApkSign")
	publishApk = GetPublishApk(this.product, this.productVersion, this.channel)
	signKeyFile := defaultSignKeyFile
	if this.channel.IsCustomSign && this.packageParam.SignKeyFile != "" {
		signKeyFile = this.packageParam.SignKeyFile
	}
	if err := android.ApkSign(signKeyFile, signPassWord, compileApk, publishApk, signAliasName); err != nil {
		panic(err)
	}
	return
}
