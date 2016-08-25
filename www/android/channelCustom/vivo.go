package channelCustom

import (
	// "encoding/json"
	// "github.com/astaxie/beego"
	"u9/tool/android"
	"u9/models"
	"github.com/astaxie/beego"
	"os"
	"u9/tool"
	"strings"
	"io/ioutil"
	)


func SetVivoMainfest(productAppEl *android.Element , product *models.Product ,packageParam *models.PackageParam){
		//修改QQ相关参数
	packageName := packageParam.PackageName
	ptAppElAc := productAppEl.GetNodeByPathAndAttr("activity", "android:name","com.bbk.payment.tenpay.VivoQQPayResultActivity")
	ptAppElIf := ptAppElAc.Node("intent-filter")
	var vivo string = "qwallet" + packageName
	vqq := ptAppElIf.GetNodeByPathAndAttr("data","android:scheme","qwalletcom.game79.mw.vivo")
	vqq.AddAttr("android:scheme",vivo) 

	ptAppElWx := productAppEl.GetNodeByPathAndAttr("activity", "android:name","com.bbk.payment.wxapi.WXPayEntryActivity")
	ptAppElWx.AddAttr("android:name",packageName + ".wxapi.WXPayEntryActivity") 

}


		
func SetVivoBuildId(product *models.Product, channel *models.Channel, packageParam *models.PackageParam,
	copyToPath,buildIdPath,packagePath,channelPath string) {
	var err error
	javaContent := `import com.bbk.payment.weixin.VivoWXPayEntryActivity;
import com.tencent.mm.sdk.modelbase.BaseReq;
import com.tencent.mm.sdk.modelbase.BaseResp;
import com.tencent.mm.sdk.openapi.IWXAPIEventHandler;
import android.content.Intent;
import android.os.Bundle;
import android.util.Log;

public class WXPayEntryActivity extends VivoWXPayEntryActivity implements IWXAPIEventHandler {

	private static final String TAG = "WXPayEntryActivity";

	@Override
	public void onCreate(Bundle savedInstanceState) {
		super.onCreate(savedInstanceState);
		this.finish();
	}

	@Override
	protected void onNewIntent(Intent intent) {
		super.onNewIntent(intent);
	}

	@Override
	public void onReq(BaseReq req) {
		// TODO Auto-generated method stub
		Log.d(TAG, "onReq, errCode = " + req);
		super.onReq(req);
	}

	@Override
	public void onResp(BaseResp resp) {
		// TODO Auto-generated method stub
		Log.d(TAG, "onPayFinish, errCode = " + resp.errCode+",resp.getType() = " + resp.getType());
		super.onResp(resp);
	}
	
}`
	content := "package " + packageParam.PackageName + ".wxapi"
	filePath := strings.Replace(packageParam.PackageName + ".wxapi", ".", "/", -1)
	smaliPath := packagePath +"/smali/"+filePath
	beego.Trace(smaliPath)
	javaPath := buildIdPath + "/src/"+filePath
	classesFile := buildIdPath + "/apk/smali/"+filePath + "/WXPayEntryActivity.smali"
	smaliFile := packagePath +"/smali/"+filePath + "/WXPayEntryActivity.smali"
	wxJar := channelPath+"/libammsdk.jar"
	cpWxJar := buildIdPath + "/libs/libammsdk.jar"
	vivoJar := channelPath+"/vivoUnionSDK_3.1.2.jar"
	cpVivoJar := buildIdPath + "/libs/vivoUnionSDK_3.1.2.jar"
	apkFile :=  buildIdPath + "/bin/project-release-unsigned.apk"
	unCompileApkPath := buildIdPath + "/apk"

	content = content + ";\r\n" + javaContent

	d1 := []byte(content)
	if err = os.MkdirAll(javaPath, 0777);err != nil{
		beego.Trace(err)
		panic(err)
	}
	if err := ioutil.WriteFile(javaPath+"/WXPayEntryActivity.java", d1, 0644);err !=nil{
		beego.Trace(err)
		panic(err)
	}
	
	if _,err = tool.CopyFile(wxJar,cpWxJar);err != nil {
		beego.Trace(err)
		panic(err)
	}
	if _,err = tool.CopyFile(vivoJar,cpVivoJar);err != nil {
		beego.Trace(err)
		panic(err)
	}

	if err := android.Ant(buildIdPath, "release"); err != nil {
		beego.Trace("ant release err:", err)
		beego.Trace("ant release err:", err)
		panic(err)
	}

	android.UnCompileApk(apkFile,unCompileApkPath)
	if err = os.MkdirAll(smaliPath, 0777);err != nil{
		beego.Trace(err)
		panic(err)
	}
	if _,err = tool.CopyFile(classesFile,smaliFile);err != nil {
		beego.Trace(err)
		panic(err)
	}
	
}