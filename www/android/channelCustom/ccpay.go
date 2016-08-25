package channelCustom

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"u9/tool/android"
	"u9/models"
	)


func SetCCPayMainfest(productAppEl *android.Element ,product *models.Product,packageParam *models.PackageParam){
	//获取参数
	jsonParam := new(map[string]interface{})
		if err := json.Unmarshal([]byte(packageParam.XmlParam), jsonParam); err != nil {
			beego.Error(err)
		}
	app_id := (*jsonParam)["app_id"].(string)
	ccsdk := productAppEl.GetNodeByPathAndAttr("activity", "android:name","com.lion.ccpay.app.user.QQPayActivity")
	ccsdkIf := ccsdk.Node("intent-filter")
	ccsdkData := ccsdkIf.GetNodeByPathAndAttr("data","android:scheme","qqPay102067")
	ccsdkData.AddAttr("android:scheme","qqPay"+app_id)
}

// func setBaiduBuildId(manifest *android.Manifest){

// }