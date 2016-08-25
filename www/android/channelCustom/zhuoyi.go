package channelCustom

import (
	// "encoding/json"
	// "github.com/astaxie/beego"
	"u9/tool/android"
	"u9/models"
	)


func SetZhuoYiMainfest(productAppEl *android.Element ,product *models.Product,packageParam *models.PackageParam){
//获取参数
	packageName := packageParam.PackageName

	zysdk := productAppEl.GetNodeByPathAndAttr("provider", "android:name","com.zhuoyi.system.promotion.provider.PromWebContentProvider")
	zysdk.AddAttr("android:authorities",packageName)

	zysdk1 := productAppEl.GetNodeByPathAndAttr("provider", "android:name","com.droi.account.DroiAccountProvider")
	zysdk1.AddAttr("android:authorities",packageName + ".droidatabase")
	
}

// func SetZhuoYiBuildId(manifest *android.Manifest){

// }

func SetZhuoYiRes(channelPath,productPath,packagePath string,product *models.Product,
	channel *models.Channel, packageParam  *models.PackageParam){
	// channelResPath := channelPath + "/" + strconv.Itoa(product.Id)
	// packageResPath := packagePath + "/res/drawable/"
	xml := packagePath + "/res/xml/" + "lib_droi_account_authenticator.xml"
	rootEl := android.LoadXmlFile(xml)
	resEl := rootEl.GetNodeByPath("account-authenticator")
	resEl.AddAttr("android:accountType",packageParam.PackageName)
	
}