package pay

import (
	"github.com/astaxie/beego"
	"strconv"
	"u9/api/channel/notify/pay"
	"u9/api/common"
	"u9/models"
)

func (this *PayController) ChannelPayNotify() {
	ret := ""
	defer func() {
		this.Ctx.WriteString(ret)
	}()

	msg := common.DumpCtx(this.Ctx)
	beego.Trace(msg)

	productId, _ := strconv.Atoi(this.Ctx.Input.Param(":productId"))
	if isExist := new(models.Product).Query().Filter("Id", productId).Exist(); !isExist {
		ret = "3007:Product:productId exception"
		return
	}

	channelId, _ := strconv.Atoi(this.Ctx.Input.Param(":channelId"))
	if isExist := new(models.Channel).Query().Filter("Id", channelId).Exist(); !isExist {
		ret = "3006:Channel:channelId exception"
		return
	}

	beego.Trace("channelPayNotify: callPayNotify")

	var err error
	if ret, err = channelPayNotify.CallPayNotify(channelId, productId, this.Ctx); err != nil {
	}
}

/*
  test url:
  http://192.168.0.185/api/productPayNotiyTest/?
  ProductId=1000&
  UserId=23c16f5323755132272fba79ab2e11d8&
  ChannelId=100&
  ChannelUserId=test10086002&
  ProductOrderId=game20160114142841787&
  OrderId=201601081906563&
  ChannelOrderId=201601081906563&
  ReqAmount=32&
  PayAmount=32&
  Sign=97a45d9c3a9787182ca71af4a6c27d68
*/
