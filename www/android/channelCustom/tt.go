package channelCustom

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/config"
	"fmt"
	"u9/models"
	"io/ioutil"
	)

func SetTTAssets(product *models.Product, packageParam *models.PackageParam,
	packagePath string){
	conf := new(config.IniConfig)
	iniConfigContainer, err := conf.ParseData([]byte(""))

	if err != nil {
		fmt.Println(err)
		return
	}
	jsonParam := new(map[string]interface{})
		if err := json.Unmarshal([]byte(packageParam.XmlParam), jsonParam); err != nil {
			beego.Error(err)
		}
	gameId := (*jsonParam)["TT_SDK_GAMEID"].(string)
	secrect :=  (*jsonParam)["TT_SDK_SECRECT"].(string)
	iniConfigContainer.Set("source", "TT")
	iniConfigContainer.Set("gameId", gameId)
	iniConfigContainer.Set("changeAccount_switch", "False")
	err = iniConfigContainer.SaveConfigFile(packagePath + "/assets/tt_game_sdk_opt.properties")
	if err != nil {
		fmt.Println(err)
		return
	}
	d1 := []byte(secrect)
	if err := ioutil.WriteFile(packagePath +"/assets/TTGameSDKConfig.cfg", d1, 0644);err !=nil{
		beego.Trace(err)
		panic(err)
	}
	beego.Trace("SetTTAssets  ok")
}