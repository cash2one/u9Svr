package androidPackage

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"io/ioutil"
	"strconv"
	"strings"
	"u9/models"
	"u9/tool/android"
	"u9/www/android/channelCustom"
)

/*
type StringResources struct {
	XMLName         xml.Name         `xml:"resources"`
	Package         string           `xml:"package,attr"`
	VersionCode     string           `xml:"http://schemas.android.com/apk/res/android versionCode,attr"`
	VersionName     string           `xml:"http://schemas.android.com/apk/res/android versionName,attr"`
	InstallLocation string           `xml:"http://schemas.android.com/apk/res/android installLocation,attr"`
	ResourceString  []ResourceString `xml:"string"`
}

type ResourceString struct {
	XMLName    xml.Name `xml:"string"`
	StringName string   `xml:"http://schemas.android.com/apk/res/android name,attr"`
	InnerText  string   `xml:",innerxml"`
}
*/
const (
	amName = "AndroidManifest.xml"
)

type Manifest struct {
	productPath string
	channelPath string
	packagePath string
	channelId int
	packageName string

	packageParam     *models.PackageParam
	channel          *models.Channel
	product          *models.Product
	channelSdkParams map[string]string
	cpSdkParams      map[string]string

	productRootEl *android.Element
	channelRootEl *android.Element
	productAppEl  *android.Element
	channelAppEl  *android.Element
}

func NewManifest(packageTaskId int,
	product *models.Product, productVersion *models.ProductVersion,
	packageParam *models.PackageParam, channel *models.Channel) *Manifest {
	ret := new(Manifest)
	ret.Init(packageTaskId, product, productVersion, packageParam, channel)
	return ret
}

func (this *Manifest) Init(packageTaskId int,
	product *models.Product, productVersion *models.ProductVersion,
	packageParam *models.PackageParam, channel *models.Channel) {

	apkName := GetApkName(product, productVersion)
	this.channelPath = GetChannelPath(channel)
	this.productPath = GetProductPath(product, apkName)
	this.packagePath = GetPackagePath(packageTaskId, apkName)
	this.channelId = channel.Id
	this.packageParam, this.channel,this.product = packageParam, channel ,product
	if err := json.Unmarshal([]byte(this.packageParam.XmlParam), &this.channelSdkParams); err != nil {
		panic(err)
	}
	if (this.packageParam.ExtParam!=""){
		if err := json.Unmarshal([]byte(this.packageParam.ExtParam), &this.cpSdkParams); err != nil {
			panic(err)
		}
	}
}

func (this *Manifest) Handle() (err error) {
	this.load()
	this.clear()
	this.setApp()
	this.setPack()
	this.merge()
	this.setMeta()
	switch this.channelId{
	case 104:
		channelCustom.SetCCPayMainfest(this.productAppEl, this.product, this.packageParam)
	case 128:
		channelCustom.SetZhuoYiMainfest(this.productAppEl, this.product, this.packageParam)
	case 138:
		fallthrough
	case 139:
		channelCustom.SetTencentMainfest(this.productAppEl, this.product, this.packageParam)
	case 144:
		channelCustom.SetVivoMainfest(this.productAppEl, this.product, this.packageParam)
	case 146:
		channelCustom.SetLenovoMainfest(this.productAppEl, this.product, this.packageParam)
	case 147:
		channelCustom.SetBaiduMainfest(this.productAppEl, this.product, this.packageParam)
	}
	ioutil.WriteFile(this.packagePath+"/"+amName, []byte(this.productRootEl.SyncToXml()), 0666)
	return
}

func (this *Manifest) merge() {
	this.productAppEl.Parent().MergeByName(this.channelAppEl.Parent(), "supports-screens")
	this.productAppEl.Parent().MergeByNameAndAttr(this.channelAppEl.Parent(),
		"uses-permission", "android:name")
	this.productAppEl.MergeByAttr(this.channelAppEl, "android:name")
}

func (this *Manifest) load() {
	this.productRootEl = android.LoadXmlFile(this.productPath + "/" + amName)
	this.channelRootEl = android.LoadXmlFile(this.channelPath + "/" + amName)
	if this.productRootEl == nil || this.channelRootEl == nil {
		panic("Manifest:load:productRootEl, channelRootEl is nil.")
	}

	this.productAppEl = this.productRootEl.GetNodeByPath("manifest/application")
	this.channelAppEl = this.channelRootEl.GetNodeByPath("manifest/application")

	if this.productAppEl == nil || this.channelAppEl == nil {
		panic("Manifest:load:productAppEl,channelAppEl is nil.")
	}
}

func (this *Manifest) clear() {
	this.productAppEl.RemoveNode("activity", "android:name", "com.hy.game.demo.HyGameDemo")
	this.productAppEl.RemoveNode("activity", "android:name", "com.hy.game.demo.FloatActivity")

	this.channelAppEl.RemoveNode("activity", "android:name", "com.example.test.demo.Game_SplashActivity")
	this.channelAppEl.RemoveNode("activity", "android:name", "com.example.test.demo.MainActivity")

	this.channelAppEl.RemoveNode("meta-data", "android:name", "HY_GAME_ID")
	this.channelAppEl.RemoveNode("meta-data", "android:name", "HY_CHANNEL_CODE")
	this.channelAppEl.RemoveNode("meta-data", "android:name", "HY_CHANNEL_TYPE")
}

func (this *Manifest) setMeta() {
	//根据产品ID/渠道ID/渠道类型设置相应meta-data
	ptMetaEl := this.productAppEl.GetNodeByPathAndAttr("meta-data", "android:name", "HY_GAME_ID")
	if ptMetaEl == nil {
		// panic("Manifest:setMetaData:ptMetaEl is nil.")
			ptMetaElnew := android.NewElement("meta-data","")
			ptMetaElnew.AddAttr("android:name", "HY_GAME_ID")
			ptMetaElnew.AddAttr("android:value", strconv.Itoa(this.packageParam.ProductId))
			this.productAppEl.AddNode(ptMetaElnew)
	}else{
		ptMetaEl.AddAttr("android:value", strconv.Itoa(this.packageParam.ProductId))
	}

	clMetaEl := this.productAppEl.GetNodeByPathAndAttr("meta-data", "android:name", "HY_CHANNEL_CODE")
	if clMetaEl == nil {
		// panic("Manifest:setMetaData:clMetaEl is nil.")
			clMetaElnew := android.NewElement("meta-data","")
			clMetaElnew.AddAttr("android:name", "HY_CHANNEL_CODE")
			clMetaElnew.AddAttr("android:value", strconv.Itoa(this.packageParam.ChannelId))
			this.productAppEl.AddNode(clMetaElnew)
	}else{
		clMetaEl.AddAttr("android:value", strconv.Itoa(this.packageParam.ChannelId))
	}

	ctMetaEl := this.productAppEl.GetNodeByPathAndAttr("meta-data", "android:name", "HY_CHANNEL_TYPE")
	if ctMetaEl == nil {
		// panic("Manifest:setMetaData:ctMetaEl is nil.")
			ctMetaElnew := android.NewElement("meta-data","")
			ctMetaElnew.AddAttr("android:name", "HY_CHANNEL_TYPE")
			ctMetaElnew.AddAttr("android:value", this.channel.Type)
			this.productAppEl.AddNode(ctMetaElnew)
	}else{
		ctMetaEl.AddAttr("android:value", this.channel.Type)
	}
	//调用热云
	this.setReYun()

	//特殊渠道meta-data 不加 \0
	for k, v := range this.channelSdkParams {
		el := this.productAppEl.GetNodeByPathAndAttr("meta-data", "android:name", k)
		if el != nil {
			switch this.channelId{
			case 102://360
				fallthrough
			case 104://虫虫
				fallthrough
			case 126://乐视
				fallthrough
			case 142://安锋
				fallthrough
			case 148://安锋
				fallthrough
			case 143://全民游戏
				el.AddAttr("android:name", k)
				el.AddAttr("android:value",v)
			default:
				el.AddAttr("android:name", k)
				el.AddAttr("android:value", "\\0"+v)
			}
			// beego.Trace(k, "#", v)
		}
	}

	beego.Trace("channelParam is OK")
	beego.Trace(this.cpSdkParams)
	if(this.cpSdkParams != nil){
	for k, v := range this.cpSdkParams {
		// beego.Trace("1")
		cl := this.productAppEl.GetNodeByPathAndAttr("meta-data", "android:name", k)
		if cl != nil {
			cl.AddAttr("android:name", k)
			cl.AddAttr("android:value", "\\0"+v)
		}else{
			clnew := android.NewElement("meta-data","")
			clnew.AddAttr("android:name", k)
			clnew.AddAttr("android:value", "\\0"+v)
			this.productAppEl.AddNode(clnew)
		}
	}
	beego.Trace("extParam is OK")
}
}

func (this *Manifest) setApp() {
	this.productAppEl.MergeAttrs(this.channelAppEl)
	this.productAppEl.AddAttr("android:icon", "@drawable/ic_launcher")
}

func (this *Manifest) setPack() {

	productPackage, _ := this.productAppEl.Parent().AttrValue("package")
	channelPackage, _ := this.channelAppEl.Parent().AttrValue("package")
	packages := strings.Split(channelPackage, ".")
	if productPackage == "" || channelPackage == "" || len(packages) == 0 {
		msg := fmt.Sprintf("Manifest:setPackage:productPackage %s or channelPackage %s is empty.",
			productPackage, channelPackage)
		panic(msg)
	}

	packageName := this.packageParam.PackageName
	if packageName == "" {
		//OPPO
		if packages[len(packages)-2] == "nearme" && packages[len(packages)-1] == "gamecenter" {
			packageName = productPackage + "." + packages[len(packages)-2] + "." + packages[len(packages)-1]
		} else { //默认
			packageName = productPackage + "." + packages[len(packages)-1]

		}
	}
	beego.Trace("packageName:", packageName)
	this.productAppEl.Parent().AddAttr("package", packageName)
	this.packageName = packageName
}

func (this *Manifest) setReYun() (err error){//热云
	jsonParam := new(map[string]string)
		if err = json.Unmarshal([]byte(this.packageParam.XmlParam), jsonParam); err != nil {
			beego.Error(err)
		}
		beego.Trace(1)

	app_id := ""
	ok := false
	if app_id,ok = (*jsonParam)["TRACK_APPID"]; ok != true {
		// beego.Error("2")
		return
	}

	reMetaEl := this.productAppEl.GetNodeByPathAndAttr("meta-data", "android:name", "TRACK_APPID")
	if reMetaEl == nil {
		// panic("Manifest:setMetaData:ctMetaEl is nil.")
			reMetaElnew := android.NewElement("meta-data","")
			reMetaElnew.AddAttr("android:name", "TRACK_APPID")
			reMetaElnew.AddAttr("android:value", app_id)
			this.productAppEl.AddNode(reMetaElnew)
	}else{
		reMetaEl.AddAttr("android:value", app_id)
	}
	return 
}

// func (this *Manifest) setCCpay(){//虫虫
// 	//获取参数
// 	jsonParam := new(map[string]interface{})
// 		if err := json.Unmarshal([]byte(this.packageParam.XmlParam), jsonParam); err != nil {
// 			beego.Error(err)
// 		}
// 	app_id := (*jsonParam)["app_id"].(string)
// 	ccsdk := this.productAppEl.GetNodeByPathAndAttr("activity", "android:name","com.lion.ccpay.app.user.QQPayActivity")
// 	ccsdkIf := ccsdk.Node("intent-filter")
// 	ccsdkData := ccsdkIf.GetNodeByPathAndAttr("data","android:scheme","qqPay102067")
// 	ccsdkData.AddAttr("android:scheme","qqPay"+app_id)
// }