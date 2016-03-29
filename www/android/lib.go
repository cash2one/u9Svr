package androidPackage

import (
	"github.com/astaxie/beego"
	"os"
	"u9/models"
)

type Lib struct {
	product        *models.Product
	channel        *models.Channel
	productVersion *models.ProductVersion

	productLibPath string
	packagePath    string
}

func NewLib(packageTaskId int, product *models.Product,
	productVersion *models.ProductVersion) *Lib {
	ret := new(Lib)
	ret.product, ret.productVersion = product, productVersion

	apkName := GetApkName(product, productVersion)
	ret.packagePath = GetPackagePath(packageTaskId, apkName)
	ret.productLibPath = GetProductPath(product, apkName)
	return ret
}

func (this *Lib) Handle() {

	this.ContrastLib()
}

func (this *Lib) ContrastLib() {
	if isDirExists(this.productLibPath + "/lib/armeabi-v7a") {
		// if isDirExists(this.packagePath + "/lib/armeabi-v7a") {
		beego.Trace("armeabi-v7a the found")
		// }
	} else {
		os.RemoveAll(this.packagePath + "/lib/armeabi-v7a")
		beego.Trace("armeabi-v7a not found Del", this.packagePath+"/lib/armeabi-v7a")
	}
	if isDirExists(this.productLibPath + "/lib/arm64-v8a") {
		// if isDirExists(this.packagePath + "/lib/arm64-v8a") {
		beego.Trace("arm64-v8a the found")
		// }
	} else {
		os.RemoveAll(this.packagePath + "/lib/arm64-v8a")
		beego.Trace("arm64-v8a not found Del", this.productLibPath+"/lib/arm64-v8a")
	}
	if isDirExists(this.productLibPath + "/lib/mips") {
		// if isDirExists(this.packagePath + "/lib/mips") {
		beego.Trace("mips the found")
		// }
	} else {
		os.RemoveAll(this.packagePath + "/lib/mips")
		beego.Trace("mips not found Del", this.packagePath+"/lib/mips")
	}
	if isDirExists(this.productLibPath + "/lib/mips64") {
		// if isDirExists(this.packagePath + "/lib/mips64") {
		beego.Trace("mips64 the found")
		// }
	} else {
		os.RemoveAll(this.packagePath + "/lib/mips64")
		beego.Trace("mips64 not found Del", this.packagePath+"/lib/mips64")
	}
	if isDirExists(this.productLibPath + "/lib/x86") {
		// if isDirExists(this.packagePath + "/lib/x86") {
		beego.Trace("x86 the found")
		// }
	} else {
		os.RemoveAll(this.packagePath + "/lib/x86")
		beego.Trace("x86 not found Del", this.packagePath+"/lib/x86")
	}
	if isDirExists(this.productLibPath + "/lib/x86_64") {
		// if isDirExists(this.packagePath + "/lib/x86_64") {
		beego.Trace("x86_64 the found")
		// }
	} else {
		os.RemoveAll(this.packagePath + "/lib/x86_64")
		beego.Trace("x86_64 not found Del", this.packagePath+"/lib/x86_64")
	}
}

func isDirExists(path string) bool {
	fi, err := os.Stat(path)

	if err != nil {
		return os.IsExist(err)
	} else {
		return fi.IsDir()
	}

	panic("not reached")
}
