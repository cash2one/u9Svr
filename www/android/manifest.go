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

	packageParam     *models.PackageParam
	channel          *models.Channel
	channelSdkParams map[string]string

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

	this.packageParam, this.channel = packageParam, channel
	if err := json.Unmarshal([]byte(this.packageParam.XmlParam), &this.channelSdkParams); err != nil {
		panic(err)
	}
}

func (this *Manifest) Handle() (err error) {
	this.load()
	this.clear()
	this.setApp()
	this.setPack()
	this.merge()
	this.setMeta()

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
		panic("Manifest:setMetaData:ptMetaEl is nil.")
	}
	ptMetaEl.AddAttr("android:value", strconv.Itoa(this.packageParam.ProductId))

	clMetaEl := this.productAppEl.GetNodeByPathAndAttr("meta-data", "android:name", "HY_CHANNEL_CODE")
	if clMetaEl == nil {
		panic("Manifest:setMetaData:clMetaEl is nil.")
	}
	clMetaEl.AddAttr("android:value", strconv.Itoa(this.packageParam.ChannelId))

	ctMetaEl := this.productAppEl.GetNodeByPathAndAttr("meta-data", "android:name", "HY_CHANNEL_TYPE")
	if ctMetaEl == nil {
		panic("Manifest:setMetaData:ctMetaEl is nil.")
	}
	ctMetaEl.AddAttr("android:value", this.channel.Type)

	for k, v := range this.channelSdkParams {
		el := this.productAppEl.GetNodeByPathAndAttr("meta-data", "android:name", k)
		if el != nil {
			el.AddAttr("android:name", k)
			el.AddAttr("android:value", "\\0"+v)
			// beego.Trace(k, "#", v)
		}
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
}
