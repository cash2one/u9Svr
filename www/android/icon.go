package androidPackage

import (
	"github.com/astaxie/beego"
	"u9/models"
	"u9/tool"
)

type Icon struct {
	product        *models.Product
	channel        *models.Channel
	productVersion *models.ProductVersion
	packageParam   *models.PackageParam

	channleIconPath string
	gameIconPath    string
	packagePath     string
}

func NewIcon(packageTaskId int, product *models.Product,
	channel *models.Channel,
	productVersion *models.ProductVersion,
	packageParam *models.PackageParam) *Icon {
	ret := new(Icon)
	ret.channel, ret.product, ret.productVersion, ret.packageParam = channel, product, productVersion, packageParam

	apkName := GetApkName(product, productVersion)
	ret.packagePath = GetPackagePath(packageTaskId, apkName)
	ret.channleIconPath, ret.gameIconPath = GetPackageIconPath(channel, packageParam, productVersion)

	return ret
}

func (this *Icon) Handle() {
	this.initPath()

	this.setPackageIcon()
}

func (this *Icon) initPath() {
	tool.CreateDir(this.packagePath + "/res/drawable")
	tool.CreateDir(this.packagePath + "/res/drawable-hdpi")
	tool.CreateDir(this.packagePath + "/res/drawable-ldpi")
	tool.CreateDir(this.packagePath + "/res/drawable-mdpi")
	tool.CreateDir(this.packagePath + "/res/drawable-xhdpi")
	tool.CreateDir(this.packagePath + "/res/drawable-xxhdpi")
}

func (this *Icon) setPackageIcon() {
	beego.Trace("gameIconPath:", this.gameIconPath)
	beego.Trace("channleIconPath:", this.channleIconPath)
	if this.channleIconPath != "" {
		if err := tool.GenerateImage(this.gameIconPath, this.channleIconPath, this.packagePath+"/res/drawable/ic_launcher.png", 192, 192); err != nil {
			beego.Trace(err)
			panic(err)
		}
		tool.GenerateImage(this.gameIconPath, this.channleIconPath, this.packagePath+"/res/drawable-hdpi/ic_launcher.png", 72, 72)
		tool.GenerateImage(this.gameIconPath, this.channleIconPath, this.packagePath+"/res/drawable-ldpi/ic_launcher.png", 36, 36)
		tool.GenerateImage(this.gameIconPath, this.channleIconPath, this.packagePath+"/res/drawable-mdpi/ic_launcher.png", 48, 48)
		tool.GenerateImage(this.gameIconPath, this.channleIconPath, this.packagePath+"/res/drawable-xhdpi/ic_launcher.png", 96, 96)
		tool.GenerateImage(this.gameIconPath, this.channleIconPath, this.packagePath+"/res/drawable-xxhdpi/ic_launcher.png", 144, 144)
	} else {
		if err := tool.ScaleImageFile(this.gameIconPath, this.packagePath+"/res/drawable/ic_launcher.png", 192, 192); err != nil {
			beego.Trace(err)
			panic(err)
		}
		tool.ScaleImageFile(this.gameIconPath, this.packagePath+"/res/drawable-hdpi/ic_launcher.png", 72, 72)
		tool.ScaleImageFile(this.gameIconPath, this.packagePath+"/res/drawable-ldpi/ic_launcher.png", 36, 36)
		tool.ScaleImageFile(this.gameIconPath, this.packagePath+"/res/drawable-mdpi/ic_launcher.png", 48, 48)
		tool.ScaleImageFile(this.gameIconPath, this.packagePath+"/res/drawable-xhdpi/ic_launcher.png", 96, 96)
		tool.ScaleImageFile(this.gameIconPath, this.packagePath+"/res/drawable-xxhdpi/ic_launcher.png", 144, 144)
	}

}
