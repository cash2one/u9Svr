package androidPackage

import (
	"os"
	"u9/models"
)

const (
	comPath = "smali/com"
)

type Smali struct {
	packagePath string
}

func NewSmali(packageTaskId int, product *models.Product, productVersion *models.ProductVersion) *Smali {
	ret := new(Smali)
	apkName := GetApkName(product, productVersion)
	ret.packagePath = GetPackagePath(packageTaskId, apkName)

	return ret
}

func (this *Smali) Prepare() {
	//3、删除u9sdk demo 文件
	demoPath := this.packagePath + "/" + comPath + "/" + "hy/game/demo"
	if err := os.RemoveAll(demoPath); err != nil {
		panic(err)
	}
	hygamePath := this.packagePath + "/" + comPath + "/" + "hygame"
	if err := os.RemoveAll(hygamePath); err != nil {
		panic(err)
	}
}

func (this *Smali) Handle() {
	//6、移除渠道demo
	examplePath := this.packagePath + "/" + comPath + "/" + "example"
	if err := os.RemoveAll(examplePath); err != nil {
		panic(err)
	}
}
