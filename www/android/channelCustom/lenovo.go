package channelCustom

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"u9/tool/android"
	"u9/models"
	)


func SetLenovoMainfest(productAppEl *android.Element , product *models.Product ,packageParam *models.PackageParam){
	//获取参数
	jsonParam := new(map[string]interface{})
		if err := json.Unmarshal([]byte(packageParam.XmlParam), jsonParam); err != nil {
			beego.Error(err)
		}
	packageName := packageParam.PackageName
	//联想要求
	appid := (*jsonParam)["lenovo.open.appid"].(string)
	ptAppElRe := productAppEl.GetNodeByPathAndAttr("receiver", "android:name","com.lenovo.lsf.gamesdk.receiver.GameSdkReceiver")
	ptAppElIf := ptAppElRe.Node("intent-filter")
	action := ptAppElIf.GetNodeByPathAndAttr("action","android:name","1603291086545.app.ln")
	action.AddAttr("android:name",appid)
	category := ptAppElIf.GetNodeByPathAndAttr("category","android:name","com.game79.mw.lenovo")
	category.AddAttr("android:name",packageName)
	//联想要求
	ptAppElRe2 := productAppEl.GetNodeByPathAndAttr("receiver", "android:name","com.lenovo.lsf.gamesdk.receiver.GameSdkAndroidLReceiver")
	ptAppElIf2 := ptAppElRe2.Node("intent-filter")
	category2 := ptAppElIf2.GetNodeByPathAndAttr("category","android:name","com.game79.mw.lenovo")
	category2.AddAttr("android:name",packageName)
	//修改主Activity
	mainActivity := (*jsonParam)["MainActivity"].(string)
	ptAppElMain := productAppEl.GetNodeByPathAndAttr("activity","android:name",mainActivity)
	ptAppElMainIf := ptAppElMain.Node("intent-filter")
	main := ptAppElMainIf.GetNodeByPathAndAttr("action","android:name","android.intent.action.MAIN")
	main.AddAttr("android:name","lenovoid.MAIN")
	launcher :=  ptAppElMainIf.GetNodeByPathAndAttr("category","android:name","android.intent.category.LAUNCHER")
	launcher.AddAttr("android:name","android.intent.category.DEFAULT")
	//闪屏页横竖屏设置
	direction := product.Direction
	var orientation string 
	if (direction == 0){
		orientation = "landscape"
	}else{
		orientation = "portrait"
	}
	welcomActivity := productAppEl.GetNodeByPathAndAttr("activity","android:name","com.lenovo.lsf.gamesdk.ui.WelcomeActivity")
	welcomActivity.AddAttr("android:screenOrientation",orientation)

}

// func setLenovoBuildId(manifest *android.Manifest){

// }