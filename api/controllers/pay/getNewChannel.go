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
		this.Data["json"] = ret
		this.ServeJSON(true)
	}()

	//msg := common.DumpCtx(this.Ctx)
	//beego.Trace(msg)

	productId, _ := strconv.Atoi(this.Ctx.Input.Param(":productId"))
	channelId, _ := strconv.Atoi(this.Ctx.Input.Param(":channelId"))

	format := "getNewChannel: %v"

	var pp models.PackageParam
	if err = pp.Query().
		Filter("channelId", channelId).
		Filter("productId", productId).
		One(&pp); err != nil {

		msg := fmt.Sprintf(format, err)
		beego.Error(msg)

		ret.Code = 1
		ret.Message = fmt.Sprintf(format, channelId, productId)
		return
	}

	args := new(map[string]string)
	if err = json.Unmarshal([]byte(pp.XmlParam), args); err != nil {

		msg := fmt.Sprintf(format, err)
		beego.Error(msg)

		ret.Code = 2
		ret.Message = fmt.Sprintf(format, pp.XmlParam)
		return
	}

	ok := false
	if ret.Ext, ok = (*args)["switchPayType"]; !ok {
		ret.Code = 3
		ret.Message = "getNewChannel: switchPayType isn't exist"
		err = errors.New(ret.Message)
		return
	}
	ret.SetCode(0)

}
