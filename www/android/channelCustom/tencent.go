package channelCustom

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/config"
	"fmt"
	"u9/tool/android"
	"u9/models"
	"u9/tool"
	"strconv"
	"os"
	"strings"
	"io/ioutil"
	)


func SetTencentMainfest(productAppEl *android.Element ,product *models.Product,packageParam *models.PackageParam){
//获取参数
	jsonParam := new(map[string]interface{})
		if err := json.Unmarshal([]byte(packageParam.XmlParam), jsonParam); err != nil {
			beego.Error(err)
		}
	//修改QQ相关参数
	ptAppElAcQQ := productAppEl.GetNodeByPathAndAttr("activity", "android:name","com.tencent.tauth.AuthActivity")
	ptAppElIfQQ := ptAppElAcQQ.Node("intent-filter")
	valueQQ := (*jsonParam)["QQ_APP_ID"].(string)
	var qq_appid string = "tencent" + valueQQ
	vqq := ptAppElIfQQ.GetNodeByPathAndAttr("data","android:scheme","tencent1105310119")
	vqq.AddAttr("android:scheme",qq_appid)
	//修改微信相关参数
	ptAppElAcWX := productAppEl.GetNodeByPathAndAttr("activity", "android:name", "com.tencent.tmgp.cqwz.wxapi.WXEntryActivity")
	ptAppElAcWX.AddAttr("android:taskAffinity",packageParam.PackageName+".diff")
	ptAppElAcWX.AddAttr("android:name",packageParam.PackageName+".wxapi.WXEntryActivity")
	ptAppElIfWX := ptAppElAcWX.Node("intent-filter")
	valueWX := (*jsonParam)["WX_APP_ID"].(string)
	beego.Trace(valueWX)
	vwx := ptAppElIfWX.GetNodeByPathAndAttr("data","android:scheme","wxa87b932b65d13d54")
	vwx.AddAttr("android:scheme",valueWX)

	mainActivity := (*jsonParam)["MainActivity"].(string)
	beego.Trace(mainActivity)
	ptAppElMain := productAppEl.GetNodeByPathAndAttr("activity","android:name",mainActivity)
	ptAppElMain.AddAttr("android:launchMode","singleTask")
	ptAppElMain.RemoveNodes("intent-filter")
}

func SetTencentBuildId(product *models.Product, channel *models.Channel, packageParam *models.PackageParam,
	copyToPath,buildIdPath,packagePath,channelPath string) {
//1、准备环境 可以直接使用BuildId init
	//2、拷贝YSDK jar包  
	//3、创建目录 
	//4、生成java 文件 编译
	//5、拷贝文件 项目\bin\classes\目录
	var err error
	content := "package " + packageParam.PackageName + ".wxapi"
	filePath := strings.Replace(packageParam.PackageName + ".wxapi", ".", "/", -1)
	smaliPath := packagePath +"/smali/"+filePath
	javaPath := buildIdPath + "/src/"+filePath
	classesFile := buildIdPath + "/apk/smali/"+filePath + "/WXEntryActivity.smali"
	smaliFile := packagePath +"/smali/"+filePath + "/WXEntryActivity.smali"
	tencentJar := channelPath+"/YSDK_Android_1.2.1_317.jar"
	cpTencetnJar := buildIdPath + "/libs/YSDK_Android_1.2.1_317.jar"
	apkFile :=  buildIdPath + "/bin/project-release-unsigned.apk"
	unCompileApkPath := buildIdPath + "/apk"
	content = content + ";\r\npublic class WXEntryActivity extends com.tencent.ysdk.module.user.impl.wx.YSDKWXEntryActivity{ }"
	d1 := []byte(content)
	if err = os.MkdirAll(javaPath, 0777);err != nil{
		beego.Trace(err)
		panic(err)
	}
	if err := ioutil.WriteFile(javaPath+"/WXEntryActivity.java", d1, 0644);err !=nil{
		beego.Trace(err)
		panic(err)
	}
	
	if _,err = tool.CopyFile(tencentJar,cpTencetnJar);err != nil {
		beego.Trace(err)
		panic(err)
	}

	if err := android.Ant(buildIdPath, "release"); err != nil {
		beego.Trace("ant release err:", err)
		beego.Trace("ant release err:", err)
		panic(err)
	}

	beego.Trace("ant ok And MikeDir:" + smaliPath)
	// os.RemoveAll(smaliPath)
	android.UnCompileApk(apkFile,unCompileApkPath)
	if err = os.MkdirAll(smaliPath, 0666);err != nil{
		beego.Trace(err)
		panic(err)
	}
	beego.Trace("classesFile:" + classesFile)
	beego.Trace("smaliFile:" + smaliFile)
	if _,err = tool.CopyFile(classesFile,smaliFile);err != nil {
		beego.Trace(err)
		panic(err)
	}

}

func SetTencentRes(channelPath,productPath,packagePath string,product *models.Product,
	channel *models.Channel, packageParam  *models.PackageParam){
	channelResPath := channelPath + "/" + strconv.Itoa(product.Id)
	packageResPath := packagePath + "/res/drawable/"
	var drawableFile  map[string]string
	var err error
	if drawableFile,err = tool.GetDirList(channelResPath,"");err != nil{
		beego.Trace(err)
		panic(err)
		}
	for fileName, filePath := range drawableFile{
       if _,err = tool.CopyFile(filePath,packageResPath+fileName);err != nil {
		beego.Trace(err)
		panic(err)
		}
    }
}

func SetTencentAssets(product *models.Product, packageParam *models.PackageParam,
	packagePath string){
	conf := new(config.IniConfig)
	iniConfigContainer, err := conf.ParseData([]byte(""))

	if err != nil {
		fmt.Println(err)
		return
	}
	xmlParam := new(map[string]interface{})
		if err := json.Unmarshal([]byte(packageParam.XmlParam), xmlParam); err != nil {
			beego.Error(err)
		}
	jsonParam := new(map[string]interface{})
		if err := json.Unmarshal([]byte(packageParam.JsonParam), jsonParam); err != nil {
			beego.Error(err)
		}
	qqAppId := (*xmlParam)["QQ_APP_ID"].(string)
	wxAppId := (*xmlParam)["WX_APP_ID"].(string)
	offerId := (*xmlParam)["OFFER_ID"].(string)
	isDebug := (*jsonParam)["tencentDebug"].(string)
	iniConfigContainer.Set("QQ_APP_ID", qqAppId)
	iniConfigContainer.Set("WX_APP_ID", wxAppId)
	iniConfigContainer.Set("OFFER_ID", offerId)
	if(isDebug == "true"){
		iniConfigContainer.Set("YSDK_URL", "https://ysdktest.qq.com")
	}else{
		iniConfigContainer.Set("YSDK_URL", "https://ysdk.qq.com")
	}
	
	err = iniConfigContainer.SaveConfigFile(packagePath+"/assets/ysdkconf.ini")
	if err != nil {
		fmt.Println(err)
		return
	}
}