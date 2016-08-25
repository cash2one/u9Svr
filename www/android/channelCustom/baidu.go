package channelCustom

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"u9/tool/android"
	"u9/models"
	)


func SetBaiduMainfest(productAppEl *android.Element ,product *models.Product,packageParam *models.PackageParam){
//获取参数
	jsonParam := new(map[string]interface{})
		if err := json.Unmarshal([]byte(packageParam.XmlParam), jsonParam); err != nil {
			beego.Error(err)
		}
	packageName := packageParam.PackageName
	//bdpsdk要求
	bdsdk := productAppEl.GetNodeByPathAndAttr("activity", "android:name","com.baidu.platformsdk.pay.channel.ali.AliPayActivity")
	bdsdkIf := bdsdk.Node("intent-filter")
	bdsdkData := bdsdkIf.GetNodeByPathAndAttr("data","android:scheme","bdpsdkcom.baidu.bdgamesdk.demo")
	bdsdkData.AddAttr("android:scheme","bdpsdk"+packageName)
	//qq支付
	qqsdk := productAppEl.GetNodeByPathAndAttr("activity", "android:name","com.baidu.platformsdk.pay.channel.qqwallet.QQPayActivity")
	qqsdkIf := qqsdk.Node("intent-filter")
	qqsdkData := qqsdkIf.GetNodeByPathAndAttr("data","android:scheme","qwalletcom.game79.mw.baidu")
	qqsdkData.AddAttr("android:scheme","qwallet"+packageName)
	//多酷SDK
	dksdk := productAppEl.GetNodeByPathAndAttr("provider", "android:name","com.duoku.platform.download.DownloadProvider")
	dksdk.AddAttr("android:authorities",packageName)
	//录屏SDK
	lpsdk := productAppEl.GetNodeByPathAndAttr("provider", "android:name","mobisocial.omlib.service.OmlibContentProvider")
	lpsdk.AddAttr("android:authorities",packageName+".provider")
	//录屏SDK
	lpsdk1 := productAppEl.GetNodeByPathAndAttr("provider", "android:name","glrecorder.Initializer")
	lpsdk1.AddAttr("android:authorities",packageName+".initializer")
}

// func setBaiduBuildId(manifest *android.Manifest){

// }