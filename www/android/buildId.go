package androidPackage

import (
	"github.com/astaxie/beego"
	"os"
	"u9/models"
	"u9/tool"
	"u9/tool/android"
	// "strings"
	// "io/ioutil"
	"u9/www/android/channelCustom"
)

const (
	projectPath = "package/project"
)

type BuildId struct {
	product        *models.Product
	productVersion *models.ProductVersion
	channel        *models.Channel
	packageParam   *models.PackageParam
	copyToPath     string
	buildIdPath    string
	packagePath    string
	channelPath    string
}

func NewBuildId(packageTaskId int,channel *models.Channel, product *models.Product,
	productVersion *models.ProductVersion , packageParam *models.PackageParam) *BuildId {
	ret := new(BuildId)
	ret.product, ret.productVersion = product, productVersion
	ret.channel = channel
	ret.packageParam = packageParam
	ret.channelPath = GetChannelPath(channel)
	apkName := GetApkName(product, productVersion)
	ret.packagePath = GetPackagePath(packageTaskId, apkName)
	ret.buildIdPath = GetBuildIdPath(packageTaskId, "/project")
	ret.copyToPath = GetBuildIdPath(packageTaskId, "")
	return ret
}

//流程：
// 初始化：
// 1、拷贝project 模板目录
// 2、拷贝打包目录 res、AndroidManifest.xml
// 反编译：

// 1、ant打包 "ant release"
// 2、反编译 classes.dex 输出至打包目录下(smali文件夹下)

func (this *BuildId) Handle() {
	switch this.channel.Id{
		case 106:
			fallthrough
		case 107:
			fallthrough
		case 122:
 			fallthrough
 		case 130:
 			fallthrough
 		case 147:
 			fallthrough
 		case 126:
 			fallthrough
 		case 151:
 			fallthrough
 		case 136:
			this.init()
			this.ant()
			this.dex()
		case 138:
			this.init()
			os.RemoveAll(this.buildIdPath+"/res/values/public.xml")
			// this.tencent()
			channelCustom.SetTencentBuildId(this.product,this.channel,this.packageParam,
				this.copyToPath,this.buildIdPath,this.packagePath,this.channelPath)
		case 139:
			this.init()
			// this.tencent()
			channelCustom.SetTencentBuildId(this.product,this.channel,this.packageParam,
				this.copyToPath,this.buildIdPath,this.packagePath,this.channelPath)
		case 144:
			this.init()
			// this.vivo()
			channelCustom.SetVivoBuildId(this.product,this.channel,this.packageParam,
				this.copyToPath,this.buildIdPath,this.packagePath,this.channelPath)
		default :
		
	}
}

func (this *BuildId) init() {
	os.RemoveAll(this.copyToPath)
	tool.CreateDir(this.copyToPath)
	tool.CopyDir(projectPath, this.copyToPath)
	tool.CopyDir(this.packagePath+"/res", this.buildIdPath)
	tool.CopyFile(this.packagePath+"/AndroidManifest.xml", this.buildIdPath+"/AndroidManifest.xml")
}
func (this *BuildId) ant() {
	beego.Trace("buildIdPath:", this.buildIdPath)
	beego.Trace("packagePath:", this.packagePath)
	if err := android.Ant(this.buildIdPath, "release"); err != nil {
		beego.Trace("ant release err:", err)
		beego.Trace("ant release err:", err)
		panic(err)
	}
}

func(this *BuildId) dex(){
		if err := android.UnCompileSmallDex(this.buildIdPath+"/bin/classes.dex", this.packagePath+"/smali"); err != nil {
		// "package/build_id/out"
		beego.Trace("UnCompileSmallDex erro:", err)
		panic(err)
	}
}

func (this *BuildId) tencent() {
	
	
}





