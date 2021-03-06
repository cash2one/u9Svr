package androidPackage

import (
	// "github.com/astaxie/beego"
	"os"
	"u9/models"
	// "u9/tool"
)

const (
	comPath = "smali/com"
	libPath = "package/lib/smali"
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
	examplePath := this.packagePath + "/" + comPath + "/" + "example/test"
	if err := os.RemoveAll(examplePath); err != nil {
		panic(err)
	}
	// if err := tool.CopyDir(libPath, this.packagePath); err != nil {
	// 	beego.Trace(err)
	// 	panic(err)
	// }
}

func (this *Smali) Channel(){

}

func (this *Smali) setTencent(){
	
}