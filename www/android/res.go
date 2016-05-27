package androidPackage

import (
	"github.com/astaxie/beego"
	"io/ioutil"
	"os"
	"strings"
	"u9/models"
	"u9/tool"
	"u9/tool/android"
)

const (
	valuesPath = "res/values"
	layoutPath = "res/layout"
)

type Res struct {
	product        *models.Product
	productVersion *models.ProductVersion
	packageParam   *models.PackageParam

	productPath string
	packagePath string
}

func NewRes(packageTaskId int,
	product *models.Product,
	productVersion *models.ProductVersion,
	packageParam *models.PackageParam) *Res {
	ret := new(Res)
	ret.product, ret.productVersion, ret.packageParam = product, productVersion, packageParam

	apkName := GetApkName(product, productVersion)
	ret.productPath = GetProductPath(product, apkName)
	ret.packagePath = GetPackagePath(packageTaskId, apkName)

	return ret
}

func (this *Res) Handle() {
	this.clear()
	this.merge()

	this.setPublicXml()
	this.setStringsXml()
}

func (this *Res) clear() {
	//8 public.xml
	publicXml := this.packagePath + "/" + valuesPath + "/" + "public.xml"
	if err := os.RemoveAll(publicXml); err != nil {
		beego.Trace(err)
		panic(err)
	}

	hygame_activityXml := this.packagePath + "/" + layoutPath +"/"+ "hygame_activity.xml"
	// beego.Trace(hygame_activityXml)
	if err := os.RemoveAll(hygame_activityXml); err != nil {
		beego.Trace(err)
		panic(err)
	}
	//9 layout目录中相关的xml
	demo_activityXml := this.packagePath + "/" + layoutPath + "/" + "hy_demo_activity.xml"
	if err := os.RemoveAll(demo_activityXml); err != nil {
		beego.Trace(err)
		panic(err)
	}
	float_activityXml := this.packagePath + "/" + layoutPath + "/" + "hy_demo_float_activity.xml"
	if err := os.RemoveAll(float_activityXml); err != nil {
		beego.Trace(err)
		panic(err)
	}
}

func (this *Res) merge() {

	productFiles, err1 := tool.GetDirList(this.productPath+"/"+valuesPath, ".xml")
	if err1 != nil {
		beego.Trace(err1)
		panic(err1)
	}

	packageFiles, err2 := tool.GetDirList(this.packagePath+"/"+valuesPath, ".xml")
	if err2 != nil {
		beego.Trace(err2)
		panic(err2)
	}

	for productFileName, productFilePath := range productFiles {
		isMerge := false
		for packageFileName, packageFilePath := range packageFiles {
			if strings.EqualFold(productFileName, packageFileName) {
				this.mergeXml(productFilePath, packageFilePath)
				isMerge = true
				break
			}
		}
		if !isMerge {
			if _, err := tool.CopyFile(productFilePath,
				this.packagePath+"/"+valuesPath+"/"+productFileName); err != nil {
				beego.Trace(err)
				panic(err)
			}
		}
	}
}

func (this *Res) mergeXml(productXml, packageXml string) {
	productRootEl := android.LoadXmlFile(productXml)
	packageRootEl := android.LoadXmlFile(packageXml)

	productEl := productRootEl.GetNodeByPath("resources")
	packageEl := packageRootEl.GetNodeByPath("resources")
	//移除产品包汇总 id.xml 中 demo 相关 id
	packageEl.RemoveNode("item", "name", "demo_title")
	packageEl.RemoveNode("item", "name", "hy_btn1")
	packageEl.RemoveNode("item", "name", "hy_btn2")
	packageEl.RemoveNode("item", "name", "hy_btn3")
	packageEl.RemoveNode("item", "name", "hy_btn4")
	packageEl.RemoveNode("item", "name", "hy_btn5")
	packageEl.RemoveNode("item", "name", "hy_btn6")
	packageEl.RemoveNode("item", "name", "hy_btn7")
	packageEl.RemoveNode("item", "name", "hy_btn8")
	packageEl.RemoveNode("item", "name", "hy_btn9")
	packageEl.RemoveNode("item", "name", "pay_name")
	packageEl.RemoveNode("item", "name", "imageView1")
	packageEl.RemoveNode("item", "name", "version_title")
	packageEl.RemoveNode("item", "name", "switch_l")
	packageEl.RemoveNode("item", "name", "switch_user")

	productEl.MergeByAttr(packageEl, "name")
	// //移除产品包汇总 String.xml 中 demo相关 String
	// packageEl.RemoveNode("string", "name", "logon_button_txt")
	// packageEl.RemoveNode("string", "name", "logoff_button_txt")
	// packageEl.RemoveNode("string", "name", "pay_button_txt")
	// packageEl.RemoveNode("string", "name", "up_role_button_txt")
	// packageEl.RemoveNode("string", "name", "exit_button_txt")
	// packageEl.RemoveNode("string", "name", "hello_world")
	// packageEl.RemoveNode("string", "name", "menu_settings")
	// packageEl.RemoveNode("string", "name", "help_tips")
	// packageEl.RemoveNode("string", "name", "levelUp")
	// packageEl.RemoveNode("string", "name", "CreateRole")
	// //移除渠道包中 string.xml 中 demo相关 string
	// packageEl.RemoveNode("string", "name", "app_name")
	productEl.RemoveNode("string", "name", "logon_button_txt")
	productEl.RemoveNode("string", "name", "logoff_button_txt")
	productEl.RemoveNode("string", "name", "pay_button_txt")
	productEl.RemoveNode("string", "name", "up_role_button_txt")
	productEl.RemoveNode("string", "name", "exit_button_txt")
	// productEl.RemoveNode("string", "name", "hello_world")
	// productEl.RemoveNode("string", "name", "menu_settings")
	productEl.RemoveNode("string", "name", "help_tips")
	productEl.RemoveNode("string", "name", "levelUp")
	productEl.RemoveNode("string", "name", "CreateRole")
	// //移除渠道包中 id.xml中的demo相关id
	productEl.RemoveNode("item", "name", "login_button")
	productEl.RemoveNode("item", "name", "pay_button")
	productEl.RemoveNode("item", "name", "up_role_button")
	productEl.RemoveNode("item", "name", "logout_button")
	productEl.RemoveNode("item", "name", "exit_button")
	productEl.RemoveNode("item", "name", "user_info_text")
	productEl.RemoveNode("item", "name", "demo_line")

	if err := ioutil.WriteFile(packageXml, []byte(productEl.SyncToXml()), 0666); err != nil {
		beego.Trace(err)
		panic(err)
	}
}

func (this *Res) setPublicXml() {
	publicXml := this.packagePath + "/" + valuesPath + "/" + "public.xml"
	rootEl := android.LoadXmlFile(publicXml)
	resEl := rootEl.GetNodeByPath("resources")

	resEl.RemoveNode("public", "name", "hy_demo_float")
	resEl.RemoveNode("public", "name", "hy_demo_activity")
	resEl.RemoveNode("public", "name", "hy_demo_float_activity")
	resEl.RemoveNode("public", "name", "hy_btn1")
	resEl.RemoveNode("public", "name", "hy_btn2")
	resEl.RemoveNode("public", "name", "hy_btn3")
	resEl.RemoveNode("public", "name", "hy_btn4")
	resEl.RemoveNode("public", "name", "hy_btn5")
	resEl.RemoveNode("public", "name", "hy_btn6")
	resEl.RemoveNode("public", "name", "hy_btn7")
	resEl.RemoveNode("public", "name", "hy_btn8")
	resEl.RemoveNode("public", "name", "hy_btn9")
	resEl.RemoveNode("public", "name", "demo_title")
	resEl.RemoveNode("public", "name", "demo_line")
	resEl.RemoveNode("public", "name", "pay_name")
	resEl.RemoveNode("public", "name", "imageView1")
	resEl.RemoveNode("public", "name", "version_title")
	resEl.RemoveNode("public", "name", "switch_l")
	resEl.RemoveNode("public", "name", "switch_user")
	resEl.RemoveNode("public", "name", "logon_button_txt")
	resEl.RemoveNode("public", "name", "logoff_button_txt")
	resEl.RemoveNode("public", "name", "pay_button_txt")
	resEl.RemoveNode("public", "name", "up_role_button_txt")
	resEl.RemoveNode("public", "name", "exit_button_txt")
	resEl.RemoveNode("public", "name", "levelUp")
	resEl.RemoveNode("public", "name", "switch_user")
	resEl.RemoveNode("public", "name", "user_info_text")
	resEl.RemoveNode("public", "name", "exit_button")
	resEl.RemoveNode("public", "name", "help_tips")
	resEl.RemoveNode("public", "name", "CreateRole")
	resEl.RemoveNode("public", "name", "pay_button")
	resEl.RemoveNode("public", "name", "up_role_button")
	resEl.RemoveNode("public", "name", "logout_button")

	if err := ioutil.WriteFile(publicXml, []byte(rootEl.SyncToXml()), 0666); err != nil {
		beego.Trace(err)
		panic(err)
	}
}

func (this *Res) setStringsXml() {
	stringXml := this.packagePath + "/" + valuesPath + "/" + "strings.xml"
	el := android.LoadXmlFile(stringXml)
	v := el.GetNodeByPathAndAttr("string", "name", "app_name")

	if len(this.packageParam.ProductName) != 0 {
		v.Value = this.packageParam.ProductName
	} else {
		v.Value = this.product.Name
	}

	if err := ioutil.WriteFile(stringXml, []byte(el.SyncToXml()), 0666); err != nil {
		beego.Trace(err)
		panic(err)
	}
}

func (this *Res) setChannel(){

}

func (this *Res) setTencent(){
	
}
