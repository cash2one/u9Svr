package pay

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"strconv"
	"u9/api/common"
	"u9/models"
)

func (this *PayController) GetNewChannel() {
	ret := new(common.BasicRet)
	var err error
	ret.Init()
	defer func() {
		if err != nil {
			beego.Error(err)
		}
		this.Data["json"] = ret
		this.ServeJSON(true)
	}()

	productId, _ := strconv.Atoi(this.Ctx.Input.Param(":productId"))
	channelId, _ := strconv.Atoi(this.Ctx.Input.Param(":channelId"))
	var pp models.PackageParam
	err = pp.Query().Filter("channelId", channelId).Filter("productId", productId).One(&pp)
	if err != nil {
		format := "packageParam exception: channelId=(%s),productId=(%s)"
		ret.Code = 1
		ret.Message = fmt.Sprintf(format, channelId, productId)
		return
	}

	args := new(map[string]string)
	if err = json.Unmarshal([]byte(pp.XmlParam), args); err != nil {
		format := "packageParam exception: xmlParam=(%s)"
		ret.Code = 2
		ret.Message = fmt.Sprintf(format, pp.XmlParam)
		return
	}

	ok := false
	if ret.Ext, ok = (*args)["switchPayType"]; !ok {
		ret.Code = 3
		ret.Message = fmt.Sprintf("packageParam switchPayType isn't exist")
		err = errors.New(ret.Message)
		return
	}
}
